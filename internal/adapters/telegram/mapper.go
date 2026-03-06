package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
	"github.com/vigmiranda/dragonki/internal/application/view"
)

func ToTelegramResponse(
	resp view.GameResponse, showMenu bool, cameFrom string,
) InputTelegram {

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

		case view.ImageBase64:
			photo = tgbotapi.FileBytes{
				Name:  resp.Image.Name,
				Bytes: utils.DecodeBase64(resp.Image.Base64),
			}
		}
	}

	var menu *tgbotapi.InlineKeyboardMarkup
	menu = mapMenu(resp.MenuKey, resp.MenuParams)

	return InputTelegram{
		Message:  resp.Text,
		Img:      photo,
		Menu:     menu,
		ShowMenu: showMenu,
		CameFrom: cameFrom,
	}
}
