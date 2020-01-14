package scdownloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bogem/id3v2"
	m3u8 "github.com/user/tgbot/convertm3u8"
)

var soundcloudAPI = "***REMOVED***"

// Download your song by provided url
func Download(songURL string) string {
	// log.Println("Looking for song id")
	clientID := soundcloudAPI
	songID := getSongID(songURL, clientID)
	// log.Println("Received song id: ", songID, ". Looking for song info")
	songInfo := getSongInfo(songID, clientID)
	// log.Println("Received song info object. Looking for song playlist url")
	songM3u8Link := getM3u8Link(songInfo, clientID)
	// log.Println("Received song playlist url. Downloading mp3")
	songTitle := songInfo["title"].(string)
	songMp3Name := getMp3(songM3u8Link, songTitle)
	// log.Println("Downloaded mp3 with name: ", songMp3Name, "Updating tags")
	updateSongTags(songMp3Name, songInfo)
	// log.Println("Updated song tags. Finishing job...")
	return songMp3Name
}

func getSongID(songURL string, clientID string) string {
	uri := fmt.Sprintf("https://api.soundcloud.com/resolve.json?url=%s&client_id=%s", songURL, clientID)
	res, err := http.Get(uri)
	if err != nil {
		log.Panic(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var songMetadata map[string]interface{}
	if err := json.Unmarshal(content, &songMetadata); err != nil {
		log.Panic("error:", err)
	}
	songID := fmt.Sprintf("%.0f", songMetadata["id"])
	return songID
}

func getSongInfo(songID string, clientID string) map[string]interface{} {
	uri := fmt.Sprintf("https://api-v2.soundcloud.com/tracks/%s?client_id=%s", songID, clientID)
	res, err := http.Get(uri)
	if err != nil {
		log.Panic(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var songInfo map[string]interface{}
	if err := json.Unmarshal(content, &songInfo); err != nil {
		log.Panic("error:", err)
	}
	return songInfo
}

func getM3u8Link(songInfo map[string]interface{}, clientID string) string {
	media := songInfo["media"].(map[string]interface{})
	transcodings := media["transcodings"].([]interface{})
	var songDlURL string
	for _, transcoding := range transcodings {
		trans := transcoding.(map[string]interface{})
		transFormat := trans["format"].(map[string]interface{})
		if transFormat["protocol"] != "hls" {
		} else if transFormat["mime_type"] == "audio/mpeg" {
			songDlURL = fmt.Sprintf("%s", trans["url"])
		}
	}
	if songDlURL == "" {
		log.Panic("Not found url for downloading!")
	}
	formattedSongDlURL := fmt.Sprintf("%s?client_id=%s", songDlURL, clientID)
	res, err := http.Get(formattedSongDlURL)
	if err != nil {
		log.Panic(err)
	}
	jsonContent, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var songURLMap map[string]string
	if err := json.Unmarshal(jsonContent, &songURLMap); err != nil {
		log.Panic("error:", err)
	}
	return string(songURLMap["url"])
}

func getMp3(link string, title string) string {
	name := fmt.Sprintf("%s.mp3", clearString(title))
	res, err := http.Get(link)
	if err != nil {
		log.Panic("[getMp3] Error: ", err)
	}
	defer res.Body.Close()
	if err := m3u8.Convert(res.Body, name); err != nil {
		log.Panic(err)
	}
	return name
}

func clearString(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`(\\|\/|\||\*|\:|\?|\"|\<|\>)`)
	result := re.ReplaceAllString(s, "-")
	return result
}

func updateSongTags(songName string, songInfo map[string]interface{}) {
	tag, err := id3v2.Open(songName, id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Panic("Error while opening mp3 file: ", err)
	}
	defer tag.Close()
	if titleObj := songInfo["title"]; titleObj != nil {
		titleStr := fmt.Sprintf("%s", titleObj)
		tag.SetTitle(titleStr)
	}
	if genreObj := songInfo["genre"]; genreObj != nil {
		genreStr := fmt.Sprintf("%s", genreObj)
		tag.SetGenre(genreStr)
	}
	if m := songInfo["publisher_metadata"].(map[string]interface{}); m != nil {
		if artistObj := m["artist"]; artistObj != nil {
			artistStr := fmt.Sprintf("%s", artistObj)
			tag.SetArtist(artistStr)
		} else if authorObj := songInfo["user"].(map[string]interface{}); authorObj != nil {
			if username := authorObj["username"]; username != nil {
				usernameStr := fmt.Sprintf("%s", username)
				tag.SetArtist(usernameStr)
			}
		}
		if albumObj := m["album_title"]; albumObj != nil {
			albumStr := fmt.Sprintf("%s", albumObj)
			tag.SetAlbum(albumStr)
		}
	}
	if rawDate := songInfo["display_date"].(string); rawDate != "" {
		dateObj, _ := time.Parse(time.RFC3339, rawDate)
		yearStr := fmt.Sprintf("%v", dateObj.Year())
		tag.SetYear(yearStr)
	}
	if artURL := songInfo["artwork_url"].(string); artURL != "" {
		res, err := http.Get(artURL)
		if err != nil {
			log.Panic(err)
		}
		artwork, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Panic(err)
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
			log.Panic(err)
		}
		hqArtwork, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Panic(err)
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
		log.Panic("Error while saving file: ", err)
	}
}
