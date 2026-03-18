package application

import "github.com/vigmiranda/dragonki/internal/domain"

func displayPlayerName(player *domain.Player, vipService *VIPService) string {
	if player == nil {
		return ""
	}
	return decoratePlayerName(int64(player.ID), player.Name, vipService)
}

func decoratePlayerName(playerID int64, name string, vipService *VIPService) string {
	if vipService == nil {
		return name
	}
	return vipService.DecoratePlayerName(playerID, name)
}
