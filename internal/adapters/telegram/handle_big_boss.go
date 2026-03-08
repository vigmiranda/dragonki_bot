package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func (h *Handler) handleBigBoss(action string, userID int64) InputTelegram {
	if !bootstrap.IsFeatureEnabled("big_boss_weekly", false) {
		return InputTelegram{
			Message:  "Big Boss semanal indisponivel no momento.",
			Img:      tgbotapi.FilePath("assets/error.png"),
			Menu:     mainMenu(),
			ShowMenu: true,
		}
	}

	if h.weeklyBossService == nil {
		return InputTelegram{
			Message:  "Modulo Big Boss indisponivel no momento.",
			Img:      tgbotapi.FilePath("assets/error.png"),
			Menu:     mainMenu(),
			ShowMenu: true,
		}
	}

	switch action {
	case "to_big_boss":
		gameResp, err := h.weeklyBossService.Open(userID)
		if err != nil {
			fmt.Printf("to_big_boss -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "big_boss_attack":
		gameResp, err := h.weeklyBossService.Attack(userID)
		if err != nil {
			fmt.Printf("big_boss_attack -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "big_boss_ranking":
		gameResp, err := h.weeklyBossService.Ranking()
		if err != nil {
			fmt.Printf("big_boss_ranking -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	}

	return InputTelegram{}
}
