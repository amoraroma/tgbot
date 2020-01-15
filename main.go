package main

import (
	"log"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/user/tgbot/scdownloader"
)

// @Okiarbot
// var tlgrmBotAPI = "***REMOVED***"

// @sc_download_bot
//var tlgrmBotAPI = "***REMOVED***"

func main() {
	var tlgrmBotAPI string
	tlgrmBotAPI, ok := os.LookupEnv("telegramAPI")
	if ok != true {
		// If there is not env with api -> use test account's api
		// @OkiarBot
		tlgrmBotAPI = "***REMOVED***"
	}
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
		if update.Message != nil {
			receivedMessage(bot, update.Message)
			continue
		} else if update.ChannelPost != nil {
			receivedMessage(bot, update.ChannelPost)
			continue
		} else {
			continue
		}
	}
}

func receivedMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatType := message.Chat.Type
	private := chatType == "private"
	var tmpMessageID int
	log.Println("Received new message", message.Text)
	chatID := message.Chat.ID
	if message.IsCommand() {
		log.Println("its command. Responding...")
		response := checkForCommands(message)
		sendMessage(bot, chatID, response)
		return
	}
	rawURL, ok := getSCLink(message.Text)
	if !ok {
		if private {
			log.Println("looks like there is no soundcloud url")
			sendMessage(bot, chatID, "Please send me a message with valid SoundCloud url or type /help for more info")
		}
		return
	}
	if private {
		tmpMessageID = sendMessage(bot, chatID, "Please wait, i'm downloading this song...")
	}
	songFile := scdownloader.Download(rawURL)
	log.Println("Received song name. Uploading to user")
	if private {
		tmpMessageID = sendMessage(bot, chatID, "Everything done. Uploading song to you...", tmpMessageID)
	}
	// Inform user about uploading
	if _, err := bot.Send(tgbotapi.NewChatAction(chatID, "upload_audio")); err != nil {
		log.Panic(err)
	}
	audioU := tgbotapi.NewAudioUpload(chatID, songFile)
	if _, err := bot.Send(audioU); err != nil {
		log.Panic(err)
	}
	if private {
		msgToDelete := tgbotapi.NewDeleteMessage(chatID, tmpMessageID)
		if _, err := bot.DeleteMessage(msgToDelete); err != nil {
			log.Panic(err)
		}
	}
	deleteFile(songFile)
	return
}

func checkForCommands(message *tgbotapi.Message) (response string) {
	switch message.Command() {
	case "help":
		response = "Send me an url and i will respond to you with attached audio.\nIf something went wrong - try again or contact with developer (link in description)"
	case "start":
		response = "Hello, #{username}.\nI'm SoundCloud downloader bot.\nSend me an url and i will respond with attached audio file"
	default:
		response = "I don't know that command. Just send me an url to SoundCloud song"
	}
	return response
}

func getSCLink(message string) (url string, ok bool) {
	re := regexp.MustCompile(`htt(p|ps)://soundcloud\.com/\S+/\S+`)
	url = re.FindString(message)
	if url == "" {
		return "", false
	}
	return url, true
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string, oldMsgContainer ...int) (msgID int) {
	var msgObj tgbotapi.Chattable
	if oldMsgContainer != nil {
		// Edit old message instead of creating new
		oldMsgID := oldMsgContainer[0]
		msgObj = tgbotapi.NewEditMessageText(chatID, oldMsgID, text)
	} else {
		msgObj = tgbotapi.NewMessage(chatID, text)
	}
	sentMsg, err := bot.Send(msgObj)
	if err != nil {
		log.Panic(err)
		return 0
	}
	return sentMsg.MessageID
}

func deleteFile(name string) {
	if err := os.Remove(name); err != nil {
		log.Panic(err)
	}
	return
}
