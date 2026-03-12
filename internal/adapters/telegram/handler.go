package telegram

import (
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic no handler:", r)
			log.Println(string(debug.Stack()))
		}
	}()

	if update.CallbackQuery != nil {
		h.handleCallback(update)
		return
	}

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	nameTelegram := update.Message.From.FirstName
	userName := update.Message.From.UserName

	status, err := h.onboardingService.HasCompletedOnboarding(update.Message.Chat.ID, userID, nameTelegram, userName)
	if err != nil {
		_, _ = h.presenter.Render(newTextRender(update.Message.Chat.ID, "Erro interno"))
		return
	}
	messageText := strings.TrimSpace(update.Message.Text)
	referralNotice := ""
	if strings.HasPrefix(messageText, "/start") {
		if h.referralService != nil {
			player, playerErr := h.onboardingService.Onboarding(update.Message.Chat.ID)
			if playerErr == nil && player != nil {
				registerResult, refErr := h.referralService.RegisterByStartParam(messageText, int64(player.ID), player.OnboardingStep)
				if refErr != nil {
					log.Printf("referral register failed | chatID=%d | playerID=%d | err=%v", update.Message.Chat.ID, player.ID, refErr)
				} else if registerResult != nil && registerResult.Linked {
					referralNotice = "Convite registrado com sucesso! Ao atingir nivel " +
						formatInt(registerResult.RequiredLevel) +
						", seu convidador recebe " +
						formatInt(registerResult.RewardGold) + " Gold e " +
						formatInt(registerResult.RewardKiCrystal) + " Ki Crystal(s)."
					if registerResult.ReferrerUserID > 0 && registerResult.ReferrerUserID != update.Message.Chat.ID {
						_, _ = h.presenter.Render(newTextRender(
							registerResult.ReferrerUserID,
							"Um novo jogador entrou pelo seu link de indicacao! "+
								"Recompensa ao atingir nivel "+
								formatInt(registerResult.RequiredLevel)+".",
						))
					}
				}
			}
		}
		switch status {
		case "DONE":
			_, _ = h.presenter.Render(newMainPanelRender(update.Message.Chat.ID, "OLD_PLAYER", referralNotice))
		default:
			_, _ = h.presenter.Render(newInitialPanelRender(update.Message.Chat.ID, referralNotice))
		}

		return
	}
	if status == "CHOOSING_NAME" {
		_, err := h.handleNameInput(update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			_, _ = h.presenter.Render(newTextRender(update.Message.Chat.ID, "Erro interno"))
			return
		}
		_, _ = h.presenter.Render(newMainPanelRender(update.Message.Chat.ID, "NEW_PLAYER", ""))
	}

	h.handleCallback(update)
}

func formatInt(v int) string {
	return strconv.Itoa(v)
}
