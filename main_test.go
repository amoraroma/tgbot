package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"testing"
)

func Test_checkForCommands(t *testing.T) {
	type args struct {
		message *tgbotapi.Message
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
	}{
		{
			"Start command",
			args{
				&tgbotapi.Message{
					Text:     "/start",
					Entities: &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
				},
			},
			MsgCommandStart,
		},
		{
			"Help command",
			args{
				&tgbotapi.Message{
					Text:     "/help",
					Entities: &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}},
				},
			},
			MsgCommandHelp,
		},
		{
			"Different command",
			args{
				&tgbotapi.Message{
					Text:     "/hello",
					Entities: &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
				},
			},
			MsgCommandUnknown,
		},
		{"Not a command", args{&tgbotapi.Message{Text: "some text"}}, MsgCommandUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResponse := checkForCommands(tt.args.message); gotResponse != tt.wantResponse {
				t.Errorf("checkForCommands() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_getSCLink(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name       string
		args       args
		wantUrl    string
		wantStatus int8
	}{
		{
			"Mobile version",
			args{"some text and https://m.soundcloud.com/user/song and some more text"},
			"https://soundcloud.com/user/song",
			0,
		},
		{
			"Web version",
			args{"some text and https://soundcloud.com/user/song and some more text"},
			"https://soundcloud.com/user/song",
			0,
		},
		{
			"Playlist url",
			args{"some text and https://m.soundcloud.com/user/sets/song and some more text"},
			"",
			2,
		},
		{
			"Not a url",
			args{"some text and its all"},
			"",
			1,
		},
		{
			"Real example with playlist",
			args{"Listen to Summer 2k17 by lilly manning on #SoundCloud\n" +
				"https://soundcloud.com/user-610681392/sets/summer-2k17"},
			"",
			2,
		},
		{
			"Real example with song",
			args{"Listen to Arcane Fantasy IV by Drawn To The Sky on #SoundCloud\n" +
				"https://soundcloud.com/drawntothesky/arcane-fantasy-iv"},
			"https://soundcloud.com/drawntothesky/arcane-fantasy-iv",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrl, gotStatus := getSCLink(tt.args.message)
			if gotUrl != tt.wantUrl {
				t.Errorf("getSCLink() gotUrl = %v, want %v", gotUrl, tt.wantUrl)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("getSCLink() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
