package main

var BotPhrase = struct {
	CmdStart, CmdHelp, CmdUnknown,
	ProcessStart, ProcessUploading,
	ErrNotURL, ErrUndefined, ErrPlaylist, ErrUnsupportedFormat, ErrUnsupportedService, ErrUnavailableSong string
}{
	// Commands
	CmdStart: "Hello, #{username}.\n" +
		"I'm SoundCloud downloader bot.\n" +
		"Send me an url and i will respond with attached audio file",
	CmdHelp: "Send me an url and i will respond to you with attached audio file.\n" +
		"|===| Currently supported services: |===|\n" +
		"= soundcloud.com [only direct song urls yet]\n" +
		"\nIf something went wrong - try again or contact with developer (link in description)",
	CmdUnknown: "I don't know that command. " +
		"Please send me an url to SoundCloud song or type /help for more info",
	// Process explaining
	ProcessStart:     "Please wait...",
	ProcessUploading: "Everything done. Uploading song to you...",
	// Exceptions
	ErrNotURL:             "Please send me a message with valid url",
	ErrUndefined:          "There is some problems with this song. Please try again or contact with developer",
	ErrPlaylist:           "Sorry, but i don't work with playlists yet. Use /help for more info",
	ErrUnsupportedFormat:  "This format unsupported yet. Use /help for more info",
	ErrUnsupportedService: "This service unsupported yet. Use /help for more info",
	ErrUnavailableSong: "Can't load this song. Make sure it is available and try again.\n" +
		"Otherwise, contact with developer using link from description",
}
