package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func (h *Handler) handlePayment(action string, userID int64) InputTelegram {
	switch action {
	case "action_buy_ki":
		if !bootstrap.IsFeatureEnabled("payments", true) {
			return InputTelegram{
				Message:  "Compra de Ki Crystals temporariamente indisponivel.",
				Img:      tgbotapi.FilePath("assets/error.png"),
				Menu:     exchangeMenu(),
				ShowMenu: true,
			}
		}
		gameResp, err := h.paymentService.HandleBuyKi(userID)
		if err != nil {
			fmt.Printf("action_buy_ki -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_buy_ki")
	case "action_copy_pix_code":
		gameResp, err := h.paymentService.HandleCopyPixCode(userID)
		if err != nil {
			fmt.Printf("action_copy_pix_code -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_copy_pix_code")
	}
	return InputTelegram{}
}
