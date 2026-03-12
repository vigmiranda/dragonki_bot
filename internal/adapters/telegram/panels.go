package telegram

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
)

type telegramPresenter struct {
	delivery telegramDelivery
}

func newTelegramPresenter(delivery telegramDelivery) *telegramPresenter {
	return &telegramPresenter{delivery: delivery}
}

func (p *telegramPresenter) Render(render telegramRender) (tgbotapi.Message, error) {
	if p == nil || p.delivery == nil {
		return tgbotapi.Message{}, errTelegramDeliveryClosed
	}

	switch render.Kind {
	case telegramRenderText:
		return p.delivery.Send(tgbotapi.NewMessage(render.ChatID, render.Text))
	case telegramRenderInitialPanel:
		return p.renderInitialPanel(render.ChatID, render.NoticeText)
	case telegramRenderMainPanel:
		return p.renderMainPanel(render.ChatID, render.PanelOrigin, render.NoticeText)
	case telegramRenderEditPanel:
		return p.renderEditPanel(render.ChatID, render.MessageID, render.Input)
	case telegramRenderAckCallback:
		p.ackCallback(render.CallbackID)
		return tgbotapi.Message{}, nil
	case telegramRenderDeleteMessage:
		p.deleteMessage(render.ChatID, render.MessageID)
		return tgbotapi.Message{}, nil
	default:
		return tgbotapi.Message{}, nil
	}
}

func (p *telegramPresenter) renderInitialPanel(chatID int64, noticeText string) (tgbotapi.Message, error) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/npcs/mestre_kalunga.png"),
	)
	photo.Caption = panelNoticeText(noticeText) +
		"Bem vindo guerreiro!! \n" +
		"Eu sou o Mestre KalungA \n" +
		"Caso queira entrar nesta jornada eu irei te acompanhar!!!\n\n" +
		"Vamos?"
	photo.ParseMode = "HTML"
	photo.ReplyMarkup = onboardingMenu()
	msg, err := p.delivery.Send(photo)
	if err != nil {
		log.Println("erro ao enviar painel inicial:", err)
	}
	return msg, err
}

func (p *telegramPresenter) renderMainPanel(chatID int64, origin string, noticeText string) (tgbotapi.Message, error) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/maps/plain_lv1.png"),
	)
	photo.Caption = panelNoticeText(noticeText)
	if origin != "NEW_PLAYER" {
		photo.Caption += "Voce voltou, guerreiro.\n\n" +
			"Seu Ki se renovou.\n"
	} else {
		photo.Caption += "Bem Vindo, guerreiro (a).\n\n" +
			"Se voce derrotar um guerreiro mais poderoso, outro aparecera em seguida!!\n" +
			"Siga em frente porque quanto mais voce apanha, mais forte voce fica!!!\n"
	}

	photo.Caption += "O mundo segue em movimento e \n" +
		"os inimigos aguardam nas sombras.\n\n" +
		"Cada batalha o aproxima de um poder maior.\n" +
		"Avance. Lute. Evolua.\n\n" +
		"O destino nao espera."
	photo.ParseMode = "HTML"
	photo.ReplyMarkup = mainMenu()
	msg, err := p.delivery.Send(photo)
	if err != nil {
		log.Println("erro ao enviar painel:", err)
	}
	return msg, err
}

func (p *telegramPresenter) renderEditPanel(chatID int64, messageID int, input InputTelegram) (tgbotapi.Message, error) {
	caption := input.Message + "\n\n"
	media := utils.NewPhotoMedia(
		input.Img,
		input.Message,
	)

	media.Caption = caption
	media.ParseMode = "HTML"
	edit := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Media: *media,
	}
	if input.ShowMenu {
		edit.ReplyMarkup = input.Menu
	}
	msg, err := p.delivery.Send(edit)
	if err != nil {
		if strings.Contains(err.Error(), "message is not modified") {
			return msg, nil
		}
		log.Println("erro ao editar mensagem:", err)
	}
	return msg, err
}

func (p *telegramPresenter) ackCallback(callbackID string) {
	if callbackID == "" {
		return
	}
	_, _ = p.delivery.Request(tgbotapi.NewCallback(callbackID, ""))
}

func (p *telegramPresenter) deleteMessage(chatID int64, messageID int) {
	if messageID <= 0 {
		return
	}
	_, _ = p.delivery.Send(tgbotapi.NewDeleteMessage(chatID, messageID))
}

func panelNoticeText(noticeText string) string {
	trimmed := strings.TrimSpace(noticeText)
	if trimmed == "" {
		return ""
	}
	return trimmed + "\n\n"
}

func (h *Handler) handleNameInput(userID int64, name string) (InputTelegram, error) {
	name = strings.TrimSpace(name)

	response := InputTelegram{
		Message: "Nome definido com sucesso! Bem-vindo ao DragonKi!",
		Img:     tgbotapi.FilePath(fmt.Sprintf("assets/skins/man_skin_basic.png")),
	}

	if len(name) < 3 || len(name) > 16 {
		response.Message = "O Nome precisa ter entre 3 e 16 caracteres!!"
		response.Img = tgbotapi.FilePath(fmt.Sprintf("error.png"))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		response.Message = "O Nome aceita somente de A-Z de 0-9!!"
		response.Img = tgbotapi.FilePath(fmt.Sprintf("error.png"))
	}

	if err := h.onboardingService.SetName(userID, name); err != nil {
		log.Println("erro ao salvar nome do player:", err)
	}

	return response, nil
}
