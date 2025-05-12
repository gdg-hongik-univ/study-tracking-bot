package main

import (
	"fmt"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord"
	"os"
	"os/signal"
)

func main() {
	bot, err := discord.Bot(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	// open session
	err = bot.Open()
	if err != nil {
		panic(err)
	}

	defer bot.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
