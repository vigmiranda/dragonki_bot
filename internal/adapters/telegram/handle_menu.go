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
			Message:  "Eventos\nConteudo especial e desafios semanais.",
			Img:      tgbotapi.FilePath("assets/maps/plain_lv1.png"),
			Menu:     eventsMenu(),
			ShowMenu: true,
		}
	}
	return InputTelegram{}
}
