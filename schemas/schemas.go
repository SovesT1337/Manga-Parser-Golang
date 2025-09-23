package schemas


// Структуры для работы с Telegram API
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	Date      int    `json:"date"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type LinkPreviewOptions struct {
	// IsDisabled	bool `json:"is_disabled"`
	// url	String	Optional. URL to use for the link preview. If empty, then the first URL found in the message text will be used
	// prefer_small_media	Boolean	Optional. True, if the media in the link preview is supposed to be shrunk; ignored if the URL isn't explicitly specified or media size change isn't supported for the preview
	PreferLargeMedia	bool `json:"prefer_large_media"`
	ShowAboveText 		bool `json:"show_above_text"`	
}

type SendMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	ParseMode string `json:"parse_mode"`
	LinkPreviewOptions LinkPreviewOptions `json:"link_preview_options"`
}