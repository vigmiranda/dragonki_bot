package telegram

type telegramRenderKind string

const (
	telegramRenderText          telegramRenderKind = "text"
	telegramRenderInitialPanel  telegramRenderKind = "initial_panel"
	telegramRenderMainPanel     telegramRenderKind = "main_panel"
	telegramRenderEditPanel     telegramRenderKind = "edit_panel"
	telegramRenderAckCallback   telegramRenderKind = "ack_callback"
	telegramRenderDeleteMessage telegramRenderKind = "delete_message"
)

type telegramRender struct {
	Kind        telegramRenderKind
	ChatID      int64
	MessageID   int
	Text        string
	NoticeText  string
	PanelOrigin string
	CallbackID  string
	Input       InputTelegram
}

func newTextRender(chatID int64, text string) telegramRender {
	return telegramRender{
		Kind:   telegramRenderText,
		ChatID: chatID,
		Text:   text,
	}
}

func newInitialPanelRender(chatID int64, noticeText string) telegramRender {
	return telegramRender{
		Kind:       telegramRenderInitialPanel,
		ChatID:     chatID,
		NoticeText: noticeText,
	}
}

func newMainPanelRender(chatID int64, origin string, noticeText string) telegramRender {
	return telegramRender{
		Kind:        telegramRenderMainPanel,
		ChatID:      chatID,
		NoticeText:  noticeText,
		PanelOrigin: origin,
	}
}

func newEditPanelRender(chatID int64, messageID int, input InputTelegram) telegramRender {
	return telegramRender{
		Kind:      telegramRenderEditPanel,
		ChatID:    chatID,
		MessageID: messageID,
		Input:     input,
	}
}

func newAckCallbackRender(callbackID string) telegramRender {
	return telegramRender{
		Kind:       telegramRenderAckCallback,
		CallbackID: callbackID,
	}
}

func newDeleteMessageRender(chatID int64, messageID int) telegramRender {
	return telegramRender{
		Kind:      telegramRenderDeleteMessage,
		ChatID:    chatID,
		MessageID: messageID,
	}
}
