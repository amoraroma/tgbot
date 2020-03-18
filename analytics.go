package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const StatServer = "https://bigbonus.pp.ua/api/"
const ErrorServer = "https://bigbonus.pp.ua/api/err.php"

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
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func analytics(msg *tgbotapi.Message) {
	var sMsg = makeStatMessage(msg)

	jsonStat, err := json.Marshal(sMsg)
	if err != nil {
		log.Println(err)
	}
	jsonReader := bytes.NewReader(jsonStat)

	if err = sendMsgToServer("Stat", jsonReader); err != nil {
		log.Panic(err)
	}

	return
}

func makeStatMessage(msg *tgbotapi.Message) (sMsg StatMessage) {
	var sChat = StatChat{}
	var sUser = StatUser{}
	if msg.Chat != nil {
		sChat = StatChat{
			msg.Chat.ID,
			msg.Chat.Type,
			msg.Chat.Title,
			msg.Chat.UserName,
			msg.Chat.FirstName,
			msg.Chat.LastName,
		}
	}
	if msg.From != nil {
		sUser = StatUser{
			msg.From.ID,
			msg.From.FirstName,
			msg.From.LastName,
			msg.From.UserName,
			msg.From.LanguageCode,
		}
	}

	sMsg = StatMessage{sUser, msg.Date, sChat, msg.Text}
	return sMsg
}

func sendMsgToServer(mode string, r io.Reader) (err error) {
	var server string
	if mode == "Error" {
		server = ErrorServer
	} else {
		server = StatServer
	}
	var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0"
	var serverResponse ServerResponse

	req, err := http.NewRequest("POST", server, r)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if err = json.Unmarshal(body, &serverResponse); err != nil {
		return
	}
	if serverResponse.Error {
		err = fmt.Errorf("%s", serverResponse.Message)
		return
	}
	return nil
}

func reportError(e error) {
	reader := strings.NewReader(e.Error())
	if err := sendMsgToServer("Error", reader); err != nil {
		log.Println(err)
	}
}
