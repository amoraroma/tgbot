package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/user/tgbot/scdownloader"
)

const (
	MsgNotSoundCloud   string = "Please send me a message with valid SoundCloud url or type /help for more info"
	MsgError           string = "Sorry, but there is some problems with this song. Please try another one or contact with developer"
	MsgPlaylist        string = "Sorry, but i don't work with playlists yet. Type /help for more info"
	MsgDownloadingSong string = "Please wait, i'm downloading this song..."
	MsgUploadingToUser string = "Everything done. Uploading song to you..."
	MsgCommandHelp     string = "Send me an url and i will respond to you with attached audio.\n" +
		"Playlists not supported yet\n\n" +
		"If something went wrong - try again or contact with developer (link in description)"
	MsgCommandStart string = "Hello, #{username}.\n" +
		"I'm SoundCloud downloader bot.\n" +
		"Send me an url and i will respond with attached audio file"
	MsgCommandUnknown string = "I don't know that command. " +
		"Please send me an url to SoundCloud song or type /help for more info"
)

var cfg Config

func main() {

	cfg = loadConfig("config.yml")
	if (cfg == Config{}) {
		log.Fatal("Can't load config")
	}
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if err := receivedMessage(bot, update.Message); err != nil {
				pingError(err, bot, update.Message.Chat.ID)
			}
			continue
		} else if update.ChannelPost != nil {
			if err := receivedMessage(bot, update.ChannelPost); err != nil {
				reportError(err)
			}
			continue
		} else {
			continue
		}
	}
}

func pingError(e error, bot *tgbotapi.BotAPI, chatID int64) {
	sendMessage(bot, chatID, MsgError)
	reportError(e)
	return
}

func receivedMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	analytics(message)
	chatType := message.Chat.Type
	private := chatType == "private"
	var tmpMessageID int
	//log.Println("Received new message", message.Text)
	chatID := message.Chat.ID
	if message.IsCommand() {
		log.Println("Received command -", message.Text, ".Responding...")
		response := checkForCommands(message)
		sendMessage(bot, chatID, response)
		return nil
	}
	rawURL, status := getSCLink(message.Text)
	if status != 0 {
		if private {
			if status == 1 {
				log.Println("Received message without link in private chat. Responding...")
				sendMessage(bot, chatID, MsgNotSoundCloud)
			} else if status == 2 {
				log.Println("Received message with playlist url in private chat. Responding...")
				sendMessage(bot, chatID, MsgPlaylist)
			}
		}
		return nil
	}
	if private {
		tmpMessageID = sendMessage(bot, chatID, MsgDownloadingSong)
	}
	log.Println("Received message with soundcloud url. Downloading song...")
	songFile, err := scdownloader.Download(rawURL, cfg.SoundCloud.Token)
	if err != nil {
		log.Printf("There is error while downloading: %s\n", err)
		return err
	}
	log.Println("Downloaded song. Uploading to user...")
	if private {
		tmpMessageID = sendMessage(bot, chatID, MsgUploadingToUser, tmpMessageID)
	}
	// Inform user about uploading
	if _, err := bot.Send(tgbotapi.NewChatAction(chatID, "upload_audio")); err != nil {
		return err
	}
	audioU := tgbotapi.NewAudioUpload(chatID, songFile)
	if _, err := bot.Send(audioU); err != nil {
		return err
	}
	if private {
		msgToDelete := tgbotapi.NewDeleteMessage(chatID, tmpMessageID)
		if _, err := bot.DeleteMessage(msgToDelete); err != nil {
			return err
		}
	}
	log.Println("Deleting file", songFile, "...")
	deleteFile(songFile)
	log.Println("Waiting for another message ~_~")
	return nil
}

func checkForCommands(message *tgbotapi.Message) (response string) {
	switch message.Command() {
	case "help":
		response = MsgCommandHelp
	case "start":
		response = MsgCommandStart
	default:
		response = MsgCommandUnknown
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
