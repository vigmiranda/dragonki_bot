package telegram

import "fmt"

func (h *Handler) handleReferral(action string, userID int64) InputTelegram {
	switch action {
	case "open_referral":
		if h.referralService == nil {
			break
		}
		gameResp, err := h.referralService.OpenReferral(userID)
		if err != nil {
			fmt.Printf("open_referral -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	}
	return InputTelegram{}
}
