package handler

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type VoiceHandler interface {
	HandlerFunc() func(*discordgo.Session, *discordgo.VoiceStateUpdate)
}

var (
	mu   sync.RWMutex
	pool = make([]VoiceHandler, 0)
)

func Register(handler VoiceHandler) {
	mu.Lock()
	defer mu.Unlock()
	pool = append(pool, handler)
}

func All() []VoiceHandler {
	mu.RLock()
	defer mu.RUnlock()
	return pool
}
