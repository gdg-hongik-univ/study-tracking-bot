package discord

import (
	"context"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/command"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/handler"
	"log"

	"github.com/bwmarrin/discordgo"
	_ "github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/command/test"
	_ "github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/handler/voice"
)

// -------- 핵심 구조체 --------
type Bot struct {
	*discordgo.Session
	guildID string // 개발·테스트용 길드 (빈 문자열이면 전역)
	ready   context.CancelFunc
}

// NewBot 세션만 만든다 (아직 핸들러 없음)
func NewBot(token, guildID string) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Bot{Session: s, guildID: guildID}, nil
}

// Start Intents → 핸들러 → Open 순으로 한 번에
func (b *Bot) Start(ctx context.Context) error {
	// 1) Intents
	b.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildVoiceStates

	// 2) 핸들러 부착
	b.attachCommands()
	b.addHandlers()

	// 3) WebSocket 열기
	if err := b.Open(); err != nil {
		return err
	}

	// 4) ctx 취소 시 Close
	go func() {
		<-ctx.Done()
		_ = b.Close()
	}()

	return nil
}

// ---------- 내부 메서드 ----------

// attachHandlers 모든 핸들러를 여기서 등록
func (b *Bot) attachCommands() {
	// (1) Ready → Slash 커맨드 등록
	b.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		raw := make([]*discordgo.ApplicationCommand, 0, len(command.All()))
		for _, c := range command.All() { // registry 패턴
			raw = append(raw, c.Build())
		}
		if _, err := s.ApplicationCommandBulkOverwrite(
			s.State.User.ID, b.guildID, raw,
		); err != nil {
			log.Fatalf("Slash 등록 실패: %v", err)
		}
		log.Println("Slash 명령어 등록 완료")
	})

	// (2) Interaction 실행 핸들러
	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		if cmd, ok := command.Lookup(i.ApplicationCommandData().Name); ok {
			_ = cmd.Run(s, i)
		}
	})
}

func (b *Bot) addHandlers() {
	handlers := handler.All()
	for _, h := range handlers {
		b.AddHandler(h.HandlerFunc())
	}
}
