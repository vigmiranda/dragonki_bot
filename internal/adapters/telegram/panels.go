package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) sendInitialPanel(chatID, userID int64) {
	status, _ := h.service.Status(userID)

	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/maps/plain_lv1.png"),
	)
	photo.Caption = status
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = mainMenu()

	h.bot.Send(photo)
}

func (h *Handler) editPanel(
	chatID int64,
	messageID int,
	extra string,
) {
	//status, _ := h.service.Status(userID)

	caption := extra + "\n\n" //+ status

	media := tgbotapi.NewInputMediaPhoto(
		tgbotapi.FilePath("assets/maps/plain_lv1.png"),
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
