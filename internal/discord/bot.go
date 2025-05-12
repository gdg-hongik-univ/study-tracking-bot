package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/handler"
)

func Bot(token string) (s *discordgo.Session, err error) {
	// create session
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// add handler
	s.AddHandler(handler.Hello("Open_Mind"))

	return s, nil
}
