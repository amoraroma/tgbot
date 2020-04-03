package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	segmentio "gopkg.in/segmentio/analytics-go.v3"
	"log"
	"strconv"
)

func NewAnalytics(token string) Analytics {
	client := segmentio.New(token)
	return Analytics{client}
}

type Analytics struct {
	client segmentio.Client
}

func (a Analytics) Identify(msg *tgbotapi.Message) {
	var firstName, lastName, username string
	firstName = msg.Chat.FirstName
	lastName = msg.Chat.LastName
	username = msg.Chat.UserName
	if msg.From != nil {
		firstName = msg.From.FirstName
		lastName = msg.From.LastName
		username = msg.From.UserName
	}
	userID := a.getUserID(msg)
	//noinspection GoUnhandledErrorResult
	if err := a.client.Enqueue(segmentio.Identify{
		UserId: strconv.Itoa(userID),
		Traits: segmentio.NewTraits().
			SetFirstName(firstName).
			SetLastName(lastName).
			SetUsername(username),
	}); err != nil {
		log.Println(err)
	}
	return
}

func (a Analytics) NewMessage(msg *tgbotapi.Message) {
	userID := a.getUserID(msg)
	//noinspection GoUnhandledErrorResult
	if err := a.client.Enqueue(segmentio.Track{
		UserId: strconv.Itoa(userID),
		Event:  "New message",
		Properties: segmentio.NewProperties().
			Set("text", msg.Text).
			Set("from", msg.Chat.Type),
	}); err != nil {
		log.Println(err)
	}
	return
}

func (a Analytics) NewError(msg *tgbotapi.Message, err error) {
	userID := a.getUserID(msg)
	//noinspection GoUnhandledErrorResult
	if err := a.client.Enqueue(segmentio.Track{
		UserId: strconv.Itoa(userID),
		Event:  "New error",
		Properties: segmentio.NewProperties().
			Set("error", err.Error()).
			Set("message", msg.Text),
	}); err != nil {
		log.Println(err)
	}
	return
}

func (a Analytics) getUserID(msg *tgbotapi.Message) int {
	userID := int(msg.Chat.ID)
	if msg.From != nil {
		userID = msg.From.ID
	}
	return userID
}
