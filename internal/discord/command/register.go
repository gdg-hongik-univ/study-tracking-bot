package command

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Command interface {
	Name() string
	Build() *discordgo.ApplicationCommand // Discord 올릴 raw 구조
	Run(*discordgo.Session, *discordgo.InteractionCreate) error
}

var (
	mu   sync.RWMutex
	pool = map[string]Command{}
)

func Register(c Command) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := pool[c.Name()]; dup {
		panic("duplicate command name: " + c.Name())
	}
	pool[c.Name()] = c
}

func All() []Command {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Command, 0, len(pool))
	for _, c := range pool {
		out = append(out, c)
	}
	return out
}

func Lookup(name string) (Command, bool) {
	mu.RLock()
	defer mu.RUnlock()
	c, ok := pool[name]
	return c, ok
}
