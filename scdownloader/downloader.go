package scdownloader

import (
	"encoding/json"
	"fmt"
	"github.com/bogem/id3v2"
	m3u8 "github.com/user/tgbot/convertm3u8"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// Download your song by provided url
func Download(songURL string, clientID string) (string, error) {
	if clientID == "" {
		return "", fmt.Errorf("no soundcloud token provided")
	}
	songID, err := getSongID(songURL, clientID)
	if err != nil {
		return "", err
	}
	//log.Println("[downloader] Received song id: ", songID, ". Looking for song info")
	songInfo, err := getSongInfo(songID, clientID)
	if err != nil {
		return "", err
	}
	//log.Println("[downloader] Received song info object. Looking for song playlist url")
	songM3u8Link, err := getM3u8Link(songInfo, clientID)
	if err != nil {
		return "", err
	}
	//log.Println("[downloader] Received song playlist url. Downloading mp3")
	songTitle := songInfo.Title
	songMp3Name, err := getMp3(songM3u8Link, songTitle)
	if err != nil {
		return "", err
	}
	//log.Println("[downloader] Downloaded mp3 with name: ", songMp3Name, "Updating tags")
	if err := updateSongTags(songMp3Name, songInfo); err != nil {
		return "", err
	}
	//log.Println("[downloader] Updated song tags. Finishing job...")
	return songMp3Name, nil
}

func getSongID(songURL string, clientID string) (string, error) {
	uri := fmt.Sprintf("https://api-v2.soundcloud.com/resolve?url=%s&client_id=%s", songURL, clientID)
	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var songMetadata SoundCloudMetadata2
	if err := json.Unmarshal(content, &songMetadata); err != nil {
		return "", err
	}
	songID := fmt.Sprintf("%d", songMetadata.ID)
	return songID, nil
}

func getSongInfo(songID string, clientID string) (SongInfo, error) {
	var songInfo SongInfo
	uri := fmt.Sprintf("https://api-v2.soundcloud.com/tracks/%s?client_id=%s", songID, clientID)
	res, err := http.Get(uri)
	if err != nil {
		return SongInfo{}, err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err := json.Unmarshal(content, &songInfo); err != nil {
		return SongInfo{}, err
	}
	return songInfo, nil
}

func getM3u8Link(songInfo SongInfo, clientID string) (string, error) {
	transcodings := songInfo.Media.Transcodings
	var songDlURL string
	for _, transcoding := range transcodings {
		format := transcoding.Format
		if format.Protocol != "hls" {
		} else if format.MimeType == "audio/mpeg" {
			songDlURL = transcoding.URL
		}
	}
	if songDlURL == "" {
		return "", fmt.Errorf("[downloader] not found url for downloading")
	}
	formattedSongDlURL := fmt.Sprintf("%s?client_id=%s", songDlURL, clientID)
	//log.Print(formattedSongDlURL)
	res, err := http.Get(formattedSongDlURL)
	if err != nil {
		return "", err
	}
	jsonContent, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var streamURL struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(jsonContent, &streamURL); err != nil {
		return "", err
	}
	return streamURL.URL, nil
}

func getMp3(link string, title string) (string, error) {
	name := fmt.Sprintf("%s.mp3", clearString(title))
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if err := m3u8.Convert(res.Body, name); err != nil {
		return "", err
	}
	return name, nil
}

func clearString(s string) string {
	s = strings.TrimSpace(s)
	//re := regexp.MustCompile(`(\\|\/|\||\*|\:|\?|\"|\<|\>)`)
	re := regexp.MustCompile(`([\\/|*:?"<>])`)
	result := re.ReplaceAllString(s, "-")
	return result
}

func updateSongTags(songName string, songInfo SongInfo) error {
	tag, err := id3v2.Open(songName, id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		return fmt.Errorf("[downloader] Error while opening mp3 file: %s", err)
	}
	defer tag.Close()
	if songTitle := songInfo.Title; songTitle != "" {
		tag.SetTitle(songTitle)
	}
	if songGenre := songInfo.Genre; songGenre != "" {
		tag.SetGenre(songGenre)
	}
	if publisher := songInfo.PublisherMetadata; (publisher != PublisherMetadata{}) {
		if artist := publisher.Artist; artist != "" {
			tag.SetArtist(artist)
		} else if user := songInfo.User; (user != User{}) {
			if username := user.Username; username != "" {
				tag.SetArtist(username)
			}
		}
		if album := publisher.AlbumTitle; album != "" {
			tag.SetAlbum(album)
		}
	}
	if songReleaseDate := songInfo.ReleaseDate; !songReleaseDate.IsZero() {
		yearStr := fmt.Sprintf("%v", songReleaseDate.Year())
		tag.SetYear(yearStr)
	} else if songDisplayDate := songInfo.DisplayDate; !songDisplayDate.IsZero() {
		//dateObj, _ := time.Parse(time.RFC3339, rawDate)
		yearStr := fmt.Sprintf("%v", songDisplayDate.Year())
		tag.SetYear(yearStr)
	} else if songLastModified := songInfo.LastModified; !songLastModified.IsZero() {
		yearStr := fmt.Sprintf("%v", songLastModified.Year())
		tag.SetYear(yearStr)
	}
	if artURL := songInfo.ArtworkURL; artURL != "" {
		res, err := http.Get(artURL)
		if err != nil {
			return err
		}
		artwork, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}
		pic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFileIcon,
			Description: "File icon",
			Picture:     artwork,
		}
		tag.AddAttachedPicture(pic)
		hqArtURL := strings.Replace(artURL, "-large.jpg", "-t500x500.jpg", 1)
		res, err = http.Get(hqArtURL)
		if err != nil {
			return err
		}
		hqArtwork, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}
		hqPic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     hqArtwork,
		}
		tag.AddAttachedPicture(hqPic)
	}
	if err := tag.Save(); err != nil {
		return fmt.Errorf("[downloader] Error while saving file: %s", err)
	}
	return nil
}
