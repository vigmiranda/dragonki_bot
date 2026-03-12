package telegram

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const telegramCallbackGuardTTL = 3 * time.Second

type telegramCallbackGuard struct {
	ttl     time.Duration
	mu      sync.Mutex
	entries map[string]time.Time
}

func newTelegramCallbackGuard(ttl time.Duration) *telegramCallbackGuard {
	if ttl <= 0 {
		ttl = telegramCallbackGuardTTL
	}
	return &telegramCallbackGuard{
		ttl:     ttl,
		entries: make(map[string]time.Time),
	}
}

func (g *telegramCallbackGuard) TryReserve(chatID int64, messageID int, userID int64, action string, params map[string]string) bool {
	if g == nil || !shouldGuardTelegramCallback(action) || messageID <= 0 {
		return true
	}

	now := time.Now().UTC()
	key := telegramCallbackGuardKey(chatID, messageID, userID, action, params)

	g.mu.Lock()
	defer g.mu.Unlock()

	g.pruneLocked(now)
	if reservedAt, exists := g.entries[key]; exists && now.Sub(reservedAt) < g.ttl {
		return false
	}
	g.entries[key] = now
	return true
}

func (g *telegramCallbackGuard) pruneLocked(now time.Time) {
	expiration := g.ttl * 2
	for key, reservedAt := range g.entries {
		if now.Sub(reservedAt) >= expiration {
			delete(g.entries, key)
		}
	}
}

func shouldGuardTelegramCallback(action string) bool {
	switch action {
	case "action_fight",
		"action_escape",
		"action_battle_use_item",
		"action_battle_transform",
		"action_post_battle_use_item",
		"action_post_battle_transform",
		"big_boss_attack",
		"start_ancestral_ki_run",
		"abandon_ancestral_ki_run",
		"ancestral_ki_chamber_continue",
		"chamber_use_item",
		"chamber_transform":
		return true
	default:
		return false
	}
}

func telegramCallbackGuardKey(chatID int64, messageID int, userID int64, action string, params map[string]string) string {
	if len(params) == 0 {
		return fmt.Sprintf("%d:%d:%d:%s", chatID, messageID, userID, action)
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, key+"="+params[key])
	}

	return fmt.Sprintf("%d:%d:%d:%s:%s", chatID, messageID, userID, action, strings.Join(pairs, "&"))
}

func telegramCallbackSemanticDedupeKey(chatID int64, messageID int, userID int64, action string, params map[string]string, now time.Time, ttl time.Duration) (string, time.Time) {
	if ttl <= 0 {
		ttl = telegramCallbackGuardTTL
	}

	windowStart := now.UTC().Truncate(ttl)
	rawKey := fmt.Sprintf(
		"%s|window=%d",
		telegramCallbackGuardKey(chatID, messageID, userID, action, params),
		windowStart.UnixNano(),
	)

	sum := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(sum[:]), windowStart
}
