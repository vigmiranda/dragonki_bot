package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type InputTelegram struct {
	Message  string
	Img      tgbotapi.RequestFileData
	Menu     *tgbotapi.InlineKeyboardMarkup
	ShowMenu bool
	CameFrom string
}
