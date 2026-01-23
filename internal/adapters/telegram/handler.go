package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	service *application.PlayerService
}

func NewHandler(bot *tgbotapi.BotAPI, service *application.PlayerService) *Handler {
	return &Handler{
		bot:     bot,
		service: service,
	}
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("❌ Panic no handler:", r)
				}
			}()

			if update.CallbackQuery != nil {
				h.handleCallback(update)
				return
			}

			if update.Message == nil {
				return
			}

			chatID := update.Message.Chat.ID
			userID := update.Message.From.ID
			text := update.Message.Text

			switch text {

			case "/start":
				h.sendInitialPanel(chatID, userID)

			default:
			}
		}()
	}
}

func (h *Handler) handleCallback(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	switch data {

	case "action_menu":
		h.editPanel(chatID, messageID, userID, "")

	case "action_cacar":
		h.editPanel(chatID, messageID, userID, "")

	case "action_status":
		h.editPanel(chatID, messageID, userID, "")

	case "action_transformar":
		h.editPanel(chatID, messageID, userID, "🔥 Transformação em breve")

	case "action_energia":
		h.editPanel(chatID, messageID, userID, "")
	}

	h.bot.Request(
		tgbotapi.NewCallback(update.CallbackQuery.ID, ""),
	)
}
func (h *Handler) sendInitialPanel(chatID, userID int64) {
	status, _ := h.service.Status(userID)

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("assets/maps/plain_lv1.png"))
	photo.Caption = status
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = mainMenu()
	h.bot.Send(photo)

}
func (h *Handler) editPanel(
	chatID int64,
	messageID int,
	userID int64,
	extra string,
) {
	status, _ := h.service.Status(userID)

	caption := status
	if extra != "" {
		caption = extra + "\n\n" + status
	}

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

func mainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Menu", "action_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Caçar", "action_cacar"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Status", "action_status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "action_transformar"),
			tgbotapi.NewInlineKeyboardButtonData("⚡ Energia", "action_energia"),
		),
	)
}
