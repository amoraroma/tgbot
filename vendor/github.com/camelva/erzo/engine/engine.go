package engine

import (
	"fmt"
	"github.com/camelva/erzo/loaders"
	"github.com/camelva/erzo/parsers"
	"github.com/camelva/erzo/utils"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
)

var _extractors []parsers.Extractor
var _loaders []loaders.Loader

const (
	_urlPattern = `((?:[a-z]{3,6}:\/\/)|(?:^|\s))` +
		`((?:[a-zA-Z0-9\-]+\.)+[a-z]{2,13})` +
		`([\.\?\=\&\%\/\w\-]*\b)`
)

type ErrNotURL struct{}

func (ErrNotURL) Error() string {
	return "there is no valid url"
}

type ErrUndefined struct{}

func (ErrUndefined) Error() string {
	return "undefined error"
}

// parsers errors
type ErrUnsupportedService struct {
	Service string
}

func (e ErrUnsupportedService) Error() string {
	return fmt.Sprintf("%s unsupported yet", e.Service)
}

type ErrUnsupportedType struct {
	parsers.ErrFormatNotSupported
}

type ErrCantFetchInfo struct {
	parsers.ErrCantContinue
}

// loaders errors
type ErrUnsupportedProtocol struct {
	Protocol string
}

func (ErrUnsupportedProtocol) Error() string {
	return "current loaders don't work with this protocol"
}

type ErrDownloadingError struct {
	Reason string
}

func (e ErrDownloadingError) Error() string {
	return fmt.Sprintf("can't download this song: %s", e.Reason)
}

func AddExtractor(x parsers.Extractor) {
	_extractors = append(_extractors, x)
}
func AddLoader(l loaders.Loader) {
	_loaders = append(_loaders, l)
}

type Engine struct {
	extractors   []parsers.Extractor
	loaders      []loaders.Loader
	outputFolder string
}

// New return new instance of Engine
func New(out string, truncate bool) *Engine {
	if (len(_extractors) < 1) || (len(_loaders) < 1) {
		// we need at least 1 extractor and 1 loader for work
		return nil
	}
	e := &Engine{
		extractors:   _extractors,
		loaders:      _loaders,
		outputFolder: out,
	}
	if truncate {
		e.Clean()
	}
	return e
}

// Clean current e.OutputFolder directory
func (e Engine) Clean() {
	os.RemoveAll(e.outputFolder)
	return
}

// Process your message. Return file name or one of this errors:
// ErrNotURL if there is no urls in your message
// ErrUnsupportedService if url belongs to unsupported service
// ErrUnsupportedType if service supported but certain type - not yet
// ErrCantFetchInfo if fatal error occurred while extracting info from url
// ErrUnsupportedProtocol if there is no downloader for this format
// ErrDownloadingError if fatal error occurred while downloading song
// ErrUndefined any other errors
func (e Engine) Process(s string) (string, error) {
	u, ok := extractURL(s)
	if !ok {
		return "", ErrNotURL{}
	}
	info, err := e.extractInfo(*u)
	if err != nil {
		return "", err
	}
	meta := createMetadata(info)
	title, err := e.downloadSong(info, meta)
	if err != nil {
		return "", err
	}
	return title, nil
}

func (e Engine) extractInfo(u url.URL) (*parsers.ExtractorInfo, error) {
	for _, xtr := range e.extractors {
		if !xtr.Compatible(u) {
			continue
		}
		info, err := xtr.Extract(u)
		if err != nil {
			switch err.(type) {
			case parsers.ErrFormatNotSupported:
				return nil, ErrUnsupportedType{err.(parsers.ErrFormatNotSupported)}
			case parsers.ErrCantContinue:
				return nil, ErrDownloadingError{err.Error()}
			default:
				return nil, ErrUndefined{}
			}
		}
		return info, nil
	}
	return nil, ErrUnsupportedService{Service: u.Hostname()}
}

func (e Engine) downloadSong(info *parsers.ExtractorInfo, metadata []string) (string, error) {
	if _, err := ioutil.ReadDir(e.outputFolder); err != nil {
		// outputFolder don't exist. Creating it...
		if err := os.Mkdir(e.outputFolder, 0700); err != nil {
			// can't create outPutFolder. Going to save files in root directory
			e.outputFolder = ""
		}
	}
	outPath := makeFilePath(e.outputFolder, info.Permalink)
	imageURL, err := url.Parse(info.Thumbnails["original"].URL)
	var thumbnail string
	if err == nil {
		res, err := utils.Fetch(imageURL)
		if err == nil {
			thumbnail = path.Join(e.outputFolder, imageURL.Path)
			if err := ioutil.WriteFile(thumbnail, res, 0644); err != nil {
				thumbnail = ""
			}
		}
	}
	var downloadingErr error
	for _, format := range info.Formats {
		u, err := url.Parse(format.Url)
		if err != nil {
			// invalid url, try another
			continue
		}
		for _, ldr := range e.loaders {
			if !ldr.Compatible(format) {
				// incompatible with loader, try another one
				continue
			}
			if err := ldr.Get(u, outPath); err != nil {
				// save err
				downloadingErr = err
				continue
			}
			if err := ldr.UpdateTags(outPath, metadata); err != nil {
				log.Println(err)
			}
			if len(thumbnail) > 0 {
				if err := ldr.AddThumbnail(outPath, thumbnail); err != nil {
					log.Println(err)
				}
				_ = os.Remove(thumbnail)
			}
			return outPath, nil
		}
	}
	if downloadingErr != nil {
		return "", ErrDownloadingError{Reason: downloadingErr.Error()}
	}
	return "", ErrUnsupportedProtocol{}
}

func createMetadata(info *parsers.ExtractorInfo) []string {
	metaMap := map[string]string{
		"title":        info.Title,
		"album":        info.Title,
		"genre":        info.Genre,
		"artist":       info.Uploader,
		"album_artist": info.Uploader,
		"track":        strconv.Itoa(1),
		"date":         strconv.Itoa(info.Timestamp.Year()),
	}
	var metadata = make([]string, 0, len(metaMap))
	for key, value := range metaMap {
		line := fmt.Sprintf("%s=%s", key, value)
		metadata = append(metadata, line)
	}
	return metadata
}

func makeFilePath(folder string, title string) string {
	fileName := fmt.Sprintf("%s.mp3", title)
	outPath := path.Join(folder, fileName)
	//if _, err := ioutil.ReadFile(outPath); err == nil {
	//	title = fmt.Sprintf("%s-copy", title)
	//	return makeFilePath(folder, title)
	//}
	return outPath
}

// extractURL trying to extract url from message
func extractURL(message string) (u *url.URL, ok bool) {
	re := regexp.MustCompile(_urlPattern)
	rawURL := re.FindString(message)
	if len(rawURL) < 1 {
		return nil, false
	}
	link, err := url.Parse(rawURL)
	if err != nil {
		return nil, false
	}
	return link, true
}
