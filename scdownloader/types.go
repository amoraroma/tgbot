package scdownloader

import "time"

//type StreamURL struct {
//	URL string `json:"url"`
//}

type SongMetadata struct {
	CommentCount        int         `json:"comment_count"`
	Downloadable        bool        `json:"downloadable"`
	Release             interface{} `json:"release"`
	CreatedAt           string      `json:"created_at"`
	Description         string      `json:"description"`
	OriginalContentSize int         `json:"original_content_size"`
	Title               string      `json:"title"`
	TrackType           interface{} `json:"track_type"`
	Duration            int         `json:"duration"`
	VideoURL            interface{} `json:"video_url"`
	OriginalFormat      string      `json:"original_format"`
	ArtworkURL          string      `json:"artwork_url"`
	Streamable          bool        `json:"streamable"`
	TagList             string      `json:"tag_list"`
	ReleaseMonth        interface{} `json:"release_month"`
	Genre               string      `json:"genre"`
	ReleaseDay          interface{} `json:"release_day"`
	DownloadURL         string      `json:"download_url"`
	ID                  int         `json:"id"`
	State               string      `json:"state"`
	RepostsCount        int         `json:"reposts_count"`
	LastModified        string      `json:"last_modified"`
	LabelName           interface{} `json:"label_name"`
	Commentable         bool        `json:"commentable"`
	Bpm                 interface{} `json:"bpm"`
	Policy              string      `json:"policy"`
	FavoritingsCount    int         `json:"favoritings_count"`
	Kind                string      `json:"kind"`
	PurchaseURL         interface{} `json:"purchase_url"`
	ReleaseYear         interface{} `json:"release_year"`
	KeySignature        interface{} `json:"key_signature"`
	Isrc                string      `json:"isrc"`
	Sharing             string      `json:"sharing"`
	URI                 string      `json:"uri"`
	DownloadCount       int         `json:"download_count"`
	License             string      `json:"license"`
	PurchaseTitle       interface{} `json:"purchase_title"`
	UserID              int         `json:"user_id"`
	EmbeddableBy        string      `json:"embeddable_by"`
	MonetizationModel   string      `json:"monetization_model"`
	WaveformURL         string      `json:"waveform_url"`
	Permalink           string      `json:"permalink"`
	PermalinkURL        string      `json:"permalink_url"`
	User                User        `json:"user"`
	LabelID             interface{} `json:"label_id"`
	StreamURL           string      `json:"stream_url"`
	PlaybackCount       int         `json:"playback_count"`
}
type SongInfo struct {
	CommentCount      int               `json:"comment_count"`
	FullDuration      int               `json:"full_duration"`
	Downloadable      bool              `json:"downloadable"`
	CreatedAt         time.Time         `json:"created_at"`
	Description       string            `json:"description"`
	Media             Media             `json:"media"`
	Title             string            `json:"title"`
	PublisherMetadata PublisherMetadata `json:"publisher_metadata"`
	Duration          int               `json:"duration"`
	HasDownloadsLeft  bool              `json:"has_downloads_left"`
	ArtworkURL        string            `json:"artwork_url"`
	Public            bool              `json:"public"`
	Streamable        bool              `json:"streamable"`
	TagList           string            `json:"tag_list"`
	DownloadURL       string            `json:"download_url"`
	Genre             string            `json:"genre"`
	ID                int               `json:"id"`
	RepostsCount      int               `json:"reposts_count"`
	State             string            `json:"state"`
	LabelName         string            `json:"label_name"`
	LastModified      time.Time         `json:"last_modified"`
	Commentable       bool              `json:"commentable"`
	Policy            string            `json:"policy"`
	Visuals           interface{}       `json:"visuals"`
	Kind              string            `json:"kind"`
	PurchaseURL       string            `json:"purchase_url"`
	Sharing           string            `json:"sharing"`
	URI               string            `json:"uri"`
	SecretToken       string            `json:"secret_token"`
	DownloadCount     int               `json:"download_count"`
	LikesCount        int               `json:"likes_count"`
	Urn               string            `json:"urn"`
	License           string            `json:"license"`
	PurchaseTitle     string            `json:"purchase_title"`
	DisplayDate       time.Time         `json:"display_date"`
	EmbeddableBy      string            `json:"embeddable_by"`
	ReleaseDate       time.Time         `json:"release_date"`
	UserID            int               `json:"user_id"`
	MonetizationModel string            `json:"monetization_model"`
	WaveformURL       string            `json:"waveform_url"`
	Permalink         string            `json:"permalink"`
	PermalinkURL      string            `json:"permalink_url"`
	User              User              `json:"user"`
	PlaybackCount     int               `json:"playback_count"`
}
type Format struct {
	Protocol string `json:"protocol"`
	MimeType string `json:"mime_type"`
}
type Transcodings struct {
	URL      string `json:"url"`
	Preset   string `json:"preset"`
	Duration int    `json:"duration"`
	Snipped  bool   `json:"snipped"`
	Format   Format `json:"format"`
	Quality  string `json:"quality"`
}
type Media struct {
	Transcodings []Transcodings `json:"transcodings"`
}
type PublisherMetadata struct {
	Urn           string `json:"urn"`
	Explicit      bool   `json:"explicit"`
	ContainsMusic bool   `json:"contains_music"`
	Artist        string `json:"artist"`
	ID            int    `json:"id"`
	AlbumTitle    string `json:"album_title"`
}
type User struct {
	AvatarURL    string `json:"avatar_url"`
	FirstName    string `json:"first_name"`
	FullName     string `json:"full_name"`
	ID           int    `json:"id"`
	Kind         string `json:"kind"`
	LastModified string `json:"last_modified"`
	LastName     string `json:"last_name"`
	Permalink    string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
	URI          string `json:"uri"`
	Urn          string `json:"urn"`
	Username     string `json:"username"`
	Verified     bool   `json:"verified"`
	City         string `json:"city"`
	CountryCode  string `json:"country_code"`
}
