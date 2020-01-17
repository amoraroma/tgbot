package main

import (
	"fmt"
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
	//log.Println("Received new message", message.Text)
	chatID := message.Chat.ID
	if message.IsCommand() {
		log.Println("Received command -", message.Text, ".Responding...")
		response := checkForCommands(message)
		sendMessage(bot, chatID, response)
		return
	}
	rawURL, status := getSCLink(message.Text)
	if status != 0 {
		if private {
			if status == 1 {
				log.Println("Received message without link in private chat. Responding...")
				sendMessage(bot, chatID, "Please send me a message with valid SoundCloud url or type /help for more info")
			} else if status == 2 {
				log.Println("Received message with playlist url in private chat. Responding...")
				sendMessage(bot, chatID, "Sorry, but i don't work with playlists yet. Type /help for more info")
			}
		}
		return
	}
	if private {
		tmpMessageID = sendMessage(bot, chatID, "Please wait, i'm downloading this song...")
	}
	log.Println("Received message with soundcloud url. Downloading song...")
	songFile := scdownloader.Download(rawURL)
	log.Println("Downloaded song. Uploading to user...")
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
	log.Println("Deleting file", songFile, "...")
	deleteFile(songFile)
	log.Println("Waiting for another message ~_~")
	return
}

func checkForCommands(message *tgbotapi.Message) (response string) {
	switch message.Command() {
	case "help":
		response = "Send me an url and i will respond to you with attached audio.\nPlaylists not supported yet\n\nIf something went wrong - try again or contact with developer (link in description)"
	case "start":
		response = "Hello, #{username}.\nI'm SoundCloud downloader bot.\nSend me an url and i will respond with attached audio file"
	default:
		response = "I don't know that command. Just send me an url to SoundCloud song or type /help for more info"
	}
	return response
}

// If failed - return empty string as <url>
// Status codes: [0 - ok, 1 - not SC, 2 - playlist url]
func getSCLink(message string) (url string, status int8) {
	// old = `(http.?://)(m\.)?(soundcloud.com)/(\S+)/(\S+)(/\S+)?`
	regStr := `(http[s]?://)?(m\.)?(soundcloud\.com/)([\w-+\.#;!]+/)([\w-+\.#;!]+)(/[\w-+\.#;!]+)?`
	re := regexp.MustCompile(regStr)
	// res contain array with result of regExp: [0] - full string,
	// [1] - protocol, [2] - "m." if exist, [3] - domain + /, [4] - user + /,
	// [5] - song (or "sets" if its playlist) [6] - / + playlist link
	res := re.FindStringSubmatch(message)
	if res == nil {
		return "", 1
	}
	if res[5] == "sets" {
		return "", 2
	}
	url = fmt.Sprintf("https://%s%s%s", res[3], res[4], res[5])
	log.Printf("%+v", url)
	return url, 0
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
