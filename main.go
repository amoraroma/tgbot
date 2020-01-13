package main

import (
	"log"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/user/tgbot/scdownloader"
)

// @Okiarbot
// var tlgrm_bot_api = "***REMOVED***"

// @sc_download_bot
var tlgrmBotAPI = "***REMOVED***"

func main() {
	bot, err := tgbotapi.NewBotAPI(tlgrmBotAPI)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatID := &update.Message.Chat.ID

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(*chatID, "")

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help":
				msg.Text = "Send me an url and i will respond to you with attached audio.\nIf something went wrong - try again or contact with developer (link in description)"
			case "start":
				msg.Text = "Hello, #{username}.\nI'm soundcloud downloader bot.\nSend me an url and i will respond with attached audio file"
			default:
				msg.Text = "I don't know that command. Just send me an url to soundcloud song"
			}
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else {
			log.Println("Received message. Checking if there is soundcloud link")
			rawURL, ok := getSCLink(update.Message.Text)
			if ok == false {
				msg := tgbotapi.NewMessage(*chatID, "Please send me a message with valid soundcloud url or type /help for more info")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
				continue
			}
			log.Println("Start downloader...")
			dwnldMsg := tgbotapi.NewMessage(*chatID, "Please wait, i'm downloading this song...")
			tempMessage, err := bot.Send(dwnldMsg)
			if err != nil {
				log.Panic(err)
			}
			songFile := scdownloader.Download(rawURL)
			log.Println("Uploading song to user...")
			upldMsg := tgbotapi.NewEditMessageText(*chatID, tempMessage.MessageID, "Everything done. Uploading song to you...")
			if tempMessage, err = bot.Send(upldMsg); err != nil {
				log.Panic(err)
			}
			bot.Send(tgbotapi.NewChatAction(*chatID, "upload_audio"))
			audioU := tgbotapi.NewAudioUpload(*chatID, songFile)
			bot.Send(audioU)
			log.Println("Uploaded!")
			msgToDelete := tgbotapi.NewDeleteMessage(*chatID, tempMessage.MessageID)
			bot.DeleteMessage(msgToDelete)
			log.Println("Cleaning...")
			deleteFile(songFile)
			log.Println("Everything done")
		}
	}
}

func getSCLink(message string) (string, bool) {
	re := regexp.MustCompile(`https:\/\/soundcloud\.com\/\S+\/\S+`)
	url := re.FindString(message)
	if url == "" {
		return "", false
	}
	return url, true
}

func deleteFile(name string) {
	if err := os.Remove(name); err != nil {
		log.Panic(err)
	}
}
