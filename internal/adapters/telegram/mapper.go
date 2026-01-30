package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application/view"
)

func ToTelegramResponse(
	resp view.GameResponse,
) (string, tgbotapi.RequestFileData, *tgbotapi.InlineKeyboardMarkup) {

	var photo tgbotapi.RequestFileData

	if resp.Image != nil {
		switch resp.Image.Type {

		case view.ImageAsset:
			photo = tgbotapi.FilePath(resp.Image.Path)

		case view.ImageBytes:
			photo = tgbotapi.FileBytes{
				Name:  resp.Image.Name,
				Bytes: resp.Image.Data,
			}
		}
	}

	var menu *tgbotapi.InlineKeyboardMarkup
	menu = mapMenu(resp.MenuKey, resp.MenuParams)

	return resp.Text, photo, menu
}
