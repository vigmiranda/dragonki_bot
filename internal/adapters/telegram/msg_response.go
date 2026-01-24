package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MsgResponse struct {
	Message string
	Img     string
	Menu    *tgbotapi.InlineKeyboardMarkup
}
