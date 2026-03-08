package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func (h *Handler) handleTraining(action string, userID int64) InputTelegram {
	switch action {
	case "to_training":
		if bootstrap.IsFeatureEnabled("ancestral_ki_chamber", false) {
			if h.chamberService == nil {
				return InputTelegram{
					Message:  "Camara do Ki Ancestral indisponivel no momento.",
					Img:      tgbotapi.FilePath("assets/error.png"),
					Menu:     mainMenu(),
					ShowMenu: true,
				}
			}
			gameResp, err := h.chamberService.Open(userID)
			if err != nil {
				fmt.Printf("to_training(chamber) -> %v\\n", err)
				break
			}
			return ToTelegramResponse(gameResp, true, "to_training")
		}
		gameResp, err := h.transformService.Training()
		if err != nil {
			fmt.Printf("to_training -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transform")
	case "list_transform":
		gameResp, err := h.transformService.ListTransform(userID)
		if err != nil {
			fmt.Printf("list_transform -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transform")

	case "action_transform":
		gameResp, err := h.transformService.Transform(userID)
		if err != nil {
			fmt.Printf("action_transform -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transform")

	}
	return InputTelegram{}
}
