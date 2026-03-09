package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *Handler) setPixCodeMessageID(userID int64, messageID int) {
	h.pixCodeMsgMu.Lock()
	defer h.pixCodeMsgMu.Unlock()
	h.pixCodeMessageByID[userID] = messageID
}

func (h *Handler) clearPixCodeMessage(userID int64, chatID int64) {
	h.pixCodeMsgMu.Lock()
	messageID, ok := h.pixCodeMessageByID[userID]
	if ok {
		delete(h.pixCodeMessageByID, userID)
	}
	h.pixCodeMsgMu.Unlock()

	if !ok || messageID <= 0 {
		return
	}

	_, _ = h.bot.Send(tgbotapi.NewDeleteMessage(chatID, messageID))
}
