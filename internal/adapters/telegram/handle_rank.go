package telegram

import "fmt"

func (h *Handler) handleRank(action string, userID int64) InputTelegram {
	switch action {
	case "open_rank":
		gameResp, err := h.rankService.OpenWeeklyRank()
		if err != nil {
			fmt.Printf("open_rank -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "rank_top_battles":
		gameResp, err := h.rankService.WeeklyTop("total_battles")
		if err != nil {
			fmt.Printf("rank_top_battles -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "rank_top_wins":
		gameResp, err := h.rankService.WeeklyTop("wins")
		if err != nil {
			fmt.Printf("rank_top_wins -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "rank_top_losses":
		gameResp, err := h.rankService.WeeklyTop("losses")
		if err != nil {
			fmt.Printf("rank_top_losses -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "rank_top_escape":
		gameResp, err := h.rankService.WeeklyTop("escape")
		if err != nil {
			fmt.Printf("rank_top_escape -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	}
	return InputTelegram{}
}
