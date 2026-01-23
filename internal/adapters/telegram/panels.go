package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) sendInitialPanel(chatID int64) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/skins/man_woman_select.png"),
	)
	photo.Caption = "Selecione o genero do seu personagem"
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = initialMenu()

	h.bot.Send(photo)
}
func (h *Handler) newPanel(chatID int64) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/maps/plain_lv1.png"),
	)
	photo.Caption = ""
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = mainMenu()
	h.bot.Send(photo)
}

func (h *Handler) editPanel(
	chatID int64,
	messageID int,
	msgResponse MsgResponse,
) {

	caption := msgResponse.Message + "\n\n" //+ status

	media := tgbotapi.NewInputMediaPhoto(
		tgbotapi.FilePath(msgResponse.Img),
	)
	media.Caption = caption
	media.ParseMode = "Markdown"

	edit := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Media: media,
	}

	menu := mainMenu()
	edit.ReplyMarkup = &menu

	h.bot.Send(edit)
}
