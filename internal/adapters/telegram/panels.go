package telegram

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
)

func (h *Handler) sendInitialPanel(chatID int64) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/skins/man_woman_select.png"),
	)
	photo.Caption = "Selecione o genero do seu personagem"
	photo.ParseMode = "HTML"
	photo.ReplyMarkup = choosingGenderMenu()
	h.bot.Send(photo)
}
func (h *Handler) handleNameInput(chatID int64, userID int64, name string) (MsgResponse, error) {
	name = strings.TrimSpace(name)

	response := MsgResponse{
		Message: "🔥 Nome definido com sucesso! Bem-vindo ao DragonKi!",
		Img:     tgbotapi.FilePath(fmt.Sprintf("assets/skins/man_skin_basic.png")),
	}

	if len(name) < 3 || len(name) > 16 {
		response.Message = "O Nome precisa ter entre 3 e 16 caracteres!!"
		response.Img = tgbotapi.FilePath(fmt.Sprintf("error.png"))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		response.Message = "O Nome aceita somente de A-Z de 0-9!!"
		response.Img = tgbotapi.FilePath(fmt.Sprintf("error.png"))
	}

	h.playerService.SetName(userID, name)

	return response, nil
}

func (h *Handler) newPanel(chatID int64, origin string) {
	photo := tgbotapi.NewPhoto(
		chatID,
		tgbotapi.FilePath("assets/maps/plain_lv1.png"),
	)
	if origin != "NEW_PLAYER" {
		photo.Caption = fmt.Sprintf(
			"🐉 Você voltou, guerreiro.\n\n" +
				"Seu Ki se renovou.\n",
		)
	} else {
		photo.Caption = fmt.Sprintf(
			"🐉 Bem Vindo, guerreiro (a).\n\n" +
				"Se você derrotar um guerreiro mais poderoso, outro aparecerá em seguida!!\n" +
				"Siga em frente porque quanto mais você apanha, mais forte você fica!!!\n",
		)
	}

	photo.Caption += "O mundo segue em movimento e \n" +
		"os inimigos aguardam nas sombras.\n\n" +
		"Cada batalha o aproxima de um poder maior.\n" +
		"⚔️ Avance. Lute. Evolua.\n\n" +
		"O destino não espera."
	photo.ParseMode = "HTML"
	photo.ReplyMarkup = mainMenu()
	_, err := h.bot.Send(photo)
	if err != nil {
		log.Fatal(err)
	}
}
func (h *Handler) editPanel(
	chatID int64,
	messageID int,
	msgResponse MsgResponse,
) {
	menu := mainMenu()
	caption := msgResponse.Message + "\n\n"
	media := utils.NewPhotoMedia(
		msgResponse.Img,
		msgResponse.Message,
	)

	media.Caption = caption
	media.ParseMode = "HTML"
	edit := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Media: *media,
	}

	if msgResponse.Menu != nil {
		menu = msgResponse.Menu
	}
	if msgResponse.CameFrom == "select_man" || msgResponse.CameFrom == "select_woman" {
		menu = nil
	}

	edit.ReplyMarkup = menu
	_, err := h.bot.Send(edit)
	if err != nil {
		if strings.Contains(err.Error(), "message is not modified") {
			return // ignora
		}
		log.Println("❌ erro ao editar mensagem:", err)
	}

}
