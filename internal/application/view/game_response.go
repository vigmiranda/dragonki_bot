package view

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ImageType string

const (
	ImageNone  ImageType = "NONE"
	ImageAsset ImageType = "ASSET"
	ImageBytes ImageType = "BYTES"
)

type Image struct {
	Type ImageType
	Path string
	Data []byte
	Name string
}
type GameResponse struct {
	Message    string
	Img        *Image
	Conditions bool
}

type MsgResponse struct {
	Message    string
	Img        tgbotapi.RequestFileData
	Menu       tgbotapi.InlineKeyboardMarkup
	Conditions bool
}
