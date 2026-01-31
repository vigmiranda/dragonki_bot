package telegram

import "fmt"

func (h *Handler) handlePayment(action string, userID int64) InputTelegram {
	switch action {
	case "action_buy_ki":
		gameResp, err := h.paymentService.HandleBuyKi(userID)
		if err != nil {
			fmt.Printf("action_buy_ki -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_buy_ki")
	}
	return InputTelegram{}
}
