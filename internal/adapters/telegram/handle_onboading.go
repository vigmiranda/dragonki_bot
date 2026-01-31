package telegram

import "fmt"

func (h *Handler) handleOnboarding(action string, userID int64) InputTelegram {
	switch action {
	case "select_man":
		gameResp, err := h.onboardingService.SetGender(userID, "MAN")
		if err != nil {
			fmt.Printf("select_man -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, false, "select_woman")

	case "select_woman":
		gameResp, err := h.onboardingService.SetGender(userID, "WOMAN")
		if err != nil {
			fmt.Printf("select_woman -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, false, "select_woman")
	}
	return InputTelegram{}
}
