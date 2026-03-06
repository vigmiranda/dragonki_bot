package application

import (
	"math/rand"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func pickMonsterByWeight(rng *rand.Rand, monsters []domain.Monster) *domain.Monster {
	totalWeight := 0
	for i := range monsters {
		if monsters[i].SpawnWeight > 0 {
			totalWeight += monsters[i].SpawnWeight
		}
	}
	if totalWeight == 0 {
		return nil
	}

	roll := rng.Intn(totalWeight)
	acc := 0

	for i := range monsters {
		acc += monsters[i].SpawnWeight
		if roll < acc {
			return &monsters[i]
		}
	}

	return nil
}

func randBetweenInclusive(rng *rand.Rand, min, max int) int {
	if max < min {
		min, max = max, min
	}
	if max == min {
		return min
	}
	return rng.Intn(max-min+1) + min
}
