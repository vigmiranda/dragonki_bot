package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleMenu(action string, userID int64) InputTelegram {
	switch action {
	case "action_menu":
		gameResp, err := h.menuService.Menu(userID)
		if err != nil {
			fmt.Printf("action_menu -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_menu")
	case "menu_adventure":
		gameResp, err := h.menuService.Menu(userID)
		if err != nil {
			fmt.Printf("menu_adventure -> %v\\n", err)
			break
		}
		return InputTelegram{
			Message:  gameResp.Text + "\n\nAventura\nEscolha sua proxima acao.",
			Img:      tgbotapi.FilePath(gameResp.Image.Path),
			Menu:     adventureMenu(),
			ShowMenu: true,
		}
	case "menu_character":
		return InputTelegram{
			Message:  "Personagem\nGerencie status, transformacao e descanso.",
			Img:      tgbotapi.FilePath("assets/maps/plain_lv1.png"),
			Menu:     characterMenu(),
			ShowMenu: true,
		}
	case "menu_inventory_scope":
		return InputTelegram{
			Message:  "Inventario\nAcesse seus itens e equipamentos.",
			Img:      tgbotapi.FilePath("assets/maps/plain_lv1.png"),
			Menu:     inventoryScopeMenu(),
			ShowMenu: true,
		}
	case "menu_economy":
		return InputTelegram{
			Message:  "Economia\nLojinha e cambio.",
			Img:      tgbotapi.FilePath("assets/maps/plain_lv1.png"),
			Menu:     economyMenu(),
			ShowMenu: true,
		}
	case "menu_events":
		return InputTelegram{
			Message: fmt.Sprintf(
				"\n\n <b>🎉 Bem-vindo ao Salão dos Eventos, viajante.</b> \n\n" +
					"Eu sou Miyara, guardiã das celebrações de Aethoria.\n" +
					"Aqui é onde guerreiros de todas as regiões se reúnem\n" +
					"para participar de desafios especiais, festivais e batalhas únicas.\n" +
					"Os eventos aparecem por tempo limitado,\n" +
					"e apenas os mais atentos conseguem aproveitar todas as recompensas.\n" +
					"Alguns trazem tesouros raros…\n" +
					"outros colocam até os mais fortes à prova.\n\n" +
					"Fique sempre atento às oportunidades.\n" +
					"Você nunca sabe quando algo extraordinário pode começar.\n\n" +
					"Prepare seu Ki.\n" +
					"Chame seus aliados.\n" +
					"E esteja pronto para o inesperado.\n\n" +
					"<b>Os eventos de Aethoria estão prestes a começar.</b> ✨🎊"),
			Img:      tgbotapi.FilePath("assets/npcs/miyara.png"),
			Menu:     eventsMenu(),
			ShowMenu: true,
		}
	}
	return InputTelegram{}
}
