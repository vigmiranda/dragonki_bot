package application

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/vigmiranda/dragonki/internal/ports"
)

type MsgResponse struct {
	Message string
	Img     string
}

type PlayerService struct {
	Repo ports.PlayerRepository
}

func NewPlayerService(repo ports.PlayerRepository) *PlayerService {
	return &PlayerService{Repo: repo}
}

func (s *PlayerService) SetGender(userID int64, gender string) (MsgResponse, error) {
	player, _ := s.Repo.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.Gender = gender
	player.RegenEnergy()
	s.Repo.Save(player)
	msgResponse.Message = fmt.Sprintf(
		"🧙 Status\n❤️ HP: %d\n⚔️ ATK: %d\n✨ XP: %d\n🏆 Nível: %d\n⚡ Energia: %d",
		player.HP, player.ATK, player.XP, player.Level, player.Energy,
	)
	msgResponse.Img = "assets/skin/" + strings.ToLower(gender) + "_skin_basic.png"

	return msgResponse, nil
}
func (s *PlayerService) Menu(userID int64) (MsgResponse, error) {
	player, _ := s.Repo.GetOrCreate(userID)
	var msgResponse MsgResponse

	player.RegenEnergy()
	msgResponse.Message = ""
	msgResponse.Img = "assets/maps/plain_lv1.png"
	return msgResponse, nil
}

func (s *PlayerService) Status(userID int64) (MsgResponse, error) {
	player, _ := s.Repo.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.RegenEnergy()
	s.Repo.Save(player)
	msgResponse.Message = fmt.Sprintf(
		"🧙 Status\n❤️ HP: %d\n⚔️ ATK: %d\n✨ XP: %d\n🏆 Nível: %d\n⚡ Energia: %d",
		player.HP, player.ATK, player.XP, player.Level, player.Energy,
	)
	msgResponse.Img = "assets/skin/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	if player.FormBattlesLeft > 0 {
		msgResponse.Img = "assets/skin/" + strings.ToLower(player.Gender) + "_ki_awakened.png"
	}

	return msgResponse, nil
}

func (s *PlayerService) Transform(userID int64) (MsgResponse, error) {
	player, _ := s.Repo.GetOrCreate(userID)

	var msgResponse MsgResponse

	logText := ""

	if player.Level < 3 {
		logText += fmt.Sprintf("❌ Transformação desbloqueia no nível 3")
		msgResponse.Img = "assets/skin/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	}
	if player.Energy < 3 {
		logText += fmt.Sprintf("⚡ Energia insuficiente")
		msgResponse.Img = "assets/skin/" + strings.ToLower(player.Gender) + "_skin_basic.png"
	}

	if player.Level >= 3 && player.Level <= 15 && player.Energy >= 3 {
		player.Form = "Ki Awakened"
		player.FormBattlesLeft = 3
		player.Energy -= 3
		logText += fmt.Sprintf("🔥 Você despertou o *Ki Awakened*!")
		msgResponse.Img = "assets/skin/" + strings.ToLower(player.Gender) + "_ki_awakened.png"
	}
	player.RegenEnergy()
	s.Repo.Save(player)
	msgResponse.Message = logText
	return msgResponse, nil
}

func (s *PlayerService) Energy(userID int64) (MsgResponse, error) {
	player, _ := s.Repo.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.RegenEnergy()
	s.Repo.Save(player)
	msgResponse.Message = "Energy em desenvolvimento"
	msgResponse.Img = "assets/maps/plain_lv1.png"
	return msgResponse, nil
}

func (s *PlayerService) Battle(userID int64) (MsgResponse, error) {
	player, _ := s.Repo.GetOrCreate(userID)

	var msgResponse MsgResponse

	player.RegenEnergy()

	if player.Energy <= 0 {
		msgResponse.Message = "⚡ Energia esgotada"
		return msgResponse, nil
	}

	player.Energy--
	monsterHP := rand.Intn(40) + 30
	msgResponse.Img = "assets/monster/n1/monster_1.png"
	msgResponse.Message = "🐉 Um monstro apareceu!\n"

	for monsterHP > 0 && player.HP > 0 {
		monsterHP -= player.ATK
		player.HP -= rand.Intn(10) + 5
	}

	if player.HP > 0 {
		xp := rand.Intn(20) + 10
		player.XP += xp
		msgResponse.Message += fmt.Sprintf("🏆 Vitória! +%d XP\n", xp)

		if player.XP >= player.Level*50 {
			player.Level++
			player.ATK += 5
			player.HP = 100
			msgResponse.Message += "⬆️ Subiu de nível!"
		}
		if player.FormBattlesLeft > 0 {
			player.FormBattlesLeft--
		}
	} else {
		player.HP = 100
		msgResponse.Message += "💀 Derrota! HP restaurado."
	}
	if player.FormBattlesLeft == 0 {
		player.Form = "Base"
	}
	s.Repo.Save(player)
	return msgResponse, nil
}
