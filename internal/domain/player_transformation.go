package domain

func KiAwakened(player *Player) *Player {
	atk := int(float64(player.ATK) * 0.1)
	def := int(float64(player.DEF) * 0.1)
	if atk < 1 {
		atk = 1
	}
	if def < 1 {
		def = 1
	}
	player.ATK += atk
	player.DEF += def
	return player
}
