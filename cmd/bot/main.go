package main

import (
	"context"
	"fmt"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	bot, err := discord.NewBot(os.Getenv("DISCORD_TOKEN"), "1198305978140065843")
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot running…  (Ctrl-C to quit)")
	<-ctx.Done()
	fmt.Println(""
	fmt.Println("Shutdown complete")
}
