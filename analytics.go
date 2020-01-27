package main

import (
	"bytes"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
)

const StatServer = "https://bigbonus.pp.ua"

type StatMessage struct {
	From StatUser
	Date int
	Chat StatChat
	Text string
}

type StatUser struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
	Language  string
}

type StatChat struct {
	ID        int64
	Type      string
	Title     string
	Username  string
	FirstName string
	LastName  string
}

type ServerResponse struct {
	error   bool
	message string
}

func analytics(msg *tgbotapi.Message) {
	var sChat = StatChat{}
	if msg.Chat != nil {
		sChat = StatChat{
			ID:        msg.Chat.ID,
			Type:      msg.Chat.Type,
			Title:     msg.Chat.Title,
			Username:  msg.Chat.UserName,
			FirstName: msg.Chat.FirstName,
			LastName:  msg.Chat.LastName,
		}
	}
	var sUser = StatUser{}
	if msg.From != nil {
		sUser = StatUser{
			ID:        msg.From.ID,
			FirstName: msg.From.FirstName,
			LastName:  msg.From.LastName,
			Username:  msg.From.UserName,
			Language:  msg.From.LanguageCode,
		}
	}

	var sMsg = StatMessage{
		sUser,
		msg.Date,
		sChat,
		msg.Text,
	}

	jsonStat, err := json.Marshal(sMsg)
	if err != nil {
		log.Panic(err)
	}
	jsonReader := bytes.NewReader(jsonStat)

	req, err := http.NewRequest("POST", StatServer, jsonReader)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	var serverResponse ServerResponse
	if err := json.Unmarshal(body, serverResponse); err != nil {
		log.Panic(err)
	}
	if serverResponse.error == true {
		log.Println(serverResponse.message)
		return
	}

	return
}
