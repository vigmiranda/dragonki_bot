package telegram

import (
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) sendInitialPanel(chatID int64) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/skins/man_woman_select.png"),
	)
	photo.Caption = "Selecione o genero do seu personagem"
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = choosingGenderMenu()
	h.bot.Send(photo)
}
func (h *Handler) handleNameInput(chatID int64, userID int64, name string) {
	name = strings.TrimSpace(name)

	if len(name) < 3 || len(name) > 16 {
		h.bot.Send(tgbotapi.NewMessage(chatID,
			"❌ O nome deve ter entre 3 e 16 caracteres.",
		))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		h.bot.Send(tgbotapi.NewMessage(chatID,
			"❌ Use apenas letras, números ou _",
		))
		return
	}

	h.playerService.SetName(userID, name)

	h.bot.Send(tgbotapi.NewMessage(chatID,
		"🔥 Nome definido com sucesso! Bem-vindo ao DragonKi!",
	))
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

	caption := msgResponse.Message + "\n\n"

	media := tgbotapi.NewInputMediaPhoto(
		tgbotapi.FilePath(msgResponse.Img),
	)
	media.Caption = caption
	//media.ParseMode = "Markdown"

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
