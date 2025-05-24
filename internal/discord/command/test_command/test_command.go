package test_command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/command"
	"strconv"
)

type TestCommand struct{}

func (c *TestCommand) Name() string                         { return "test" }
func (c *TestCommand) Build() *discordgo.ApplicationCommand { return commands }
func (c *TestCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	sub := data.Options[0]
	var res string

	switch sub.Name {
	case "echo":
		text := ""
		time := int64(1)
		for _, opt := range sub.Options {
			switch opt.Name {
			case "message":
				text = opt.StringValue()
			case "times":
				time = opt.IntValue()
			}
		}
		for i := int64(0); i < time; i++ {
			res += text
		}
	case "add":
		a := sub.Options[0].IntValue()
		b := sub.Options[1].IntValue()
		res = strconv.FormatInt(a+b, 10)
	case "sub":
		a, b := int64(0), int64(0)
		for _, opt := range sub.Options {
			switch opt.Name {
			case "a":
				a = opt.IntValue()
			case "b":
				b = opt.IntValue()
			}
		}
		res = strconv.FormatInt(a-b, 10)
	}

	return s.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: res,
			},
		})
}

var commands = &discordgo.ApplicationCommand{
	Name:        "test",
	Description: "test command",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "echo",
			Description: "echoing",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "message to echo",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "times",
					Description: "times to echo",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "a + b",
			Options: []*discordgo.ApplicationCommandOption{
				intOpt("a"), intOpt("b"),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "sub",
			Description: "a - b",
			Options: []*discordgo.ApplicationCommandOption{
				intOpt("a"), intOpt("b"),
			},
		},
	},
}

func intOpt(name string) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type: discordgo.ApplicationCommandOptionInteger,
		Name: name, Description: name, Required: true,
	}
}

func init() {
	command.Register(&TestCommand{})
}
