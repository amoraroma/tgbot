// Package m3u8 helps you convert your *.(m3u|m3u8) files
// into standalone mp3 file.
package m3u8

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Convert(in interface{}, out string) error {
	var err error
	fileName, isString := in.(string)
	if isString {
		err = convertFile(fileName, out)
	} else {
		fileReader, isReader := in.(io.Reader)
		if isReader {
			err = convertReader(fileReader, out)
		}
	}
	return err
}

// ConvertFile read file with provided name and push received data to Convert()
func convertFile(name, out string) error {
	playlist, err := openPlaylist(name)
	if err != nil {
		return err
	}
	if err := convertReader(playlist, out); err != nil {
		return err
	}
	return nil
}

// Convert proceed your raw data and convert it to file with <out> name
func convertReader(data io.Reader, out string) error {
	songs, err := parseM3u8(data)
	if err != nil {
		return err
	}
	songData, err := joinSongs(songs)
	if err != nil {
		return err
	}
	if filepath.Ext(out) == "" {
		out = fmt.Sprintf("%s.mp3", out)
	} else if filepath.Ext(out) != ".mp3" {
		log.Panic("[convert] Invalid output file extension")
	}
	if err := writeSong(songData, out); err != nil {
		log.Panic(err)
	}
	return nil
}

func openPlaylist(file string) (io.Reader, error) {
	var fileU8 string
	if ext := filepath.Ext(file); ext == "" {
		fileU8 = fmt.Sprintf("%s.m3u8", file)
		file = fmt.Sprintf("%s.m3u", file)
	} else if ext != ".m3u" {
		return nil, errors.New("[convert] pass a correct .m3u file")
	}
	reader, err := os.Open(file)
	if err != nil {
		readerAlt, errSecond := os.Open(fileU8)
		if errSecond != nil {
			return nil, errSecond
		}
		reader = readerAlt
	}
	return reader, nil
}

func parseM3u8(r io.Reader) ([]string, error) {
	b := bufio.NewReader(r)
	p := make([]string, 0)
	for {
		rawLine, err := b.ReadBytes('\n')
		last := false

		if err == io.EOF {
			last = true
		} else if err != nil {
			return nil, err
		}

		line := strings.TrimSpace(string(rawLine))
		length := len(line)
		if length >= 1 && line[0] == '#' {
			continue
		} else if length != 0 {
			p = append(p, line)
		}

		if last {
			break
		}
	}
	return p, nil
}

func joinSongs(songsList []string) ([]byte, error) {
	songs := make([][]byte, 0)
	for _, s := range songsList {
		songData, err := getSongData(s)
		if err != nil {
			return nil, err
		}
		songs = append(songs, songData)
	}
	completeSong := bytes.Join(songs, []byte(""))
	return completeSong, nil
}

func getSongData(songURL string) ([]byte, error) {
	var songData = make([]byte, 0)
	var err error
	if strings.Contains(songURL, "http") {
		songData, err = getRemoteSongData(songURL)
		if err != nil {
			return nil, err
		}
	} else {
		songData, err = getLocalSongData(songURL)
		if err != nil {
			return nil, err
		}
	}
	return songData, nil
}

func getRemoteSongData(songURL string) ([]byte, error) {
	res, err := http.Get(songURL)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return content, nil
}

func getLocalSongData(songURL string) ([]byte, error) {
	content, err := ioutil.ReadFile(songURL)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func writeSong(songData []byte, outName string) error {
	outFile, err := os.OpenFile(outName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if _, err := outFile.Write(songData); err != nil {
		return err
	}
	return nil
}
