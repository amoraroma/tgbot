package ffmpeg

import (
	"github.com/camelva/erzo/engine"
	"github.com/camelva/erzo/parsers"
	"net/url"
	"os"
	"os/exec"
)

type loader struct {
	name      string
	bin       string
	protocols []string
}

func init() {
	protocols := []string{"http", "https", "hls", "progressive"}
	bin := findBin()
	if len(bin) < 1 {
		// not found binary, so don't init loader
		return
	}
	config = loader{"ffmpeg", bin, protocols}
	engine.AddLoader(config)
}

func (l loader) Name() string {
	return l.name
}
func (l loader) Bin() string {
	return l.bin
}
func (l loader) Get(u *url.URL, outName string) error {
	args := make([]string, 0)
	args = append(args,
		"-y",             // overwrite existing
		"-i", u.String(), // input file
		"-vn",
		"-ar", "44100",
		"-ac", "2",
		"-b:a", "128k")
	args = append(args, outName) // output name should always be latest element
	// if need debug
	//args = append(args, "-report")
	_, err := execute(l.Bin(), args...)
	if err != nil {
		return err
	}
	return nil
}
func (l loader) UpdateTags(filename string, metadata []string) error {
	args := make([]string, 0)
	args = append(args,
		"-y",
		"-i", filename,
		"-c", "copy")
	for _, el := range metadata {
		args = append(args, "-metadata", el)
	}
	args = append(args, "temp.mp3") // output name should always be latest element
	// if need debug
	//args = append(args, "-report")
	_, err := execute(l.Bin(), args...)
	//noinspection GoUnhandledErrorResult
	defer os.Rename("temp.mp3", filename)
	if err != nil {
		return err
	}
	return nil
}
func (l loader) AddThumbnail(filename string, thumb string) error {
	tempName := "temp.mp3"
	args := make([]string, 0)
	args = append(args,
		"-y",           // overwrite existing
		"-i", filename, // input file
		"-i", thumb,
		//"-vn",
		"-id3v2_version", "3",
		"-c", "copy",
		"-map", "0", "-map", "1",
		"-metadata:s:v", "title=Album cover",
		"-metadata:s:v", "comment=Cover (front)",
		tempName)
	//args = append(args, "-report")
	_, err := execute(l.Bin(), args...)
	if err != nil {
		return err
	}
	_ = os.Rename(tempName, filename)
	return nil
}

func (l loader) Compatible(f parsers.Format) bool {
	for _, p := range l.protocols {
		if p != f.Protocol {
			continue
		}
		return true
	}
	return false
}

var config loader

func findBin() string {
	bin := "ffmpeg"
	path, err := exec.LookPath(bin)
	if err != nil {
		path = bin
	}
	return path
}

func execute(app string, args ...string) ([]byte, error) {
	cmd := exec.Command(app, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}
	return out, nil
}
