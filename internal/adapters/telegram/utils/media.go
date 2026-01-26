package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Aceita qualquer tipo de mídia suportado pelo Telegram:
// - tgbotapi.FilePath
// - tgbotapi.FileBytes
// - tgbotapi.FileID
func NewPhotoMedia(
	source tgbotapi.RequestFileData,
	caption string,
) *tgbotapi.InputMediaPhoto {

	media := tgbotapi.NewInputMediaPhoto(source)

	if caption != "" {
		media.Caption = caption
		media.ParseMode = "HTML"
	}

	return &media
}
