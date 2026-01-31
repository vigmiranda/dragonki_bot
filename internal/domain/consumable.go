package domain

type ConsumableEffect string

const (
	EffectRecoverHP     ConsumableEffect = "RECOVER_HP"
	EffectRecoverEnergy ConsumableEffect = "RECOVER_ENERGY"
	EffectBoostAttack   ConsumableEffect = "BOOST_ATTACK"
)

type Consumable struct {
	Item
	Effect ConsumableEffect
	Value  int // quanto recupera ou aumenta
}
