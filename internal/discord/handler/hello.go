package handler

import "github.com/bwmarrin/discordgo"

func Hello(name string) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(discord *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == discord.State.User.ID {
			return
		}
		if m.Content != "hello" {
			return
		}

		_, err := discord.ChannelMessageSend(m.ChannelID, "Hello "+name)
		if err != nil {
			panic(err)
		}
	}
}
