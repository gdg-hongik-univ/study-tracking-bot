package voice

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/discord/handler"
	"log"
	"sync"
	"time"
)

type joinInfo struct {
	joinedAt  time.Time
	voiceChan string
}

type TimeHandler struct {
	// 음성 채널 ID → 텍스트 채널 ID
	VoiceToText map[string]string

	// 유저별 입장 정보 저장
	mu        sync.Mutex
	joinTimes map[string]*joinInfo
}

func (t *TimeHandler) HandlerFunc() func(*discordgo.Session, *discordgo.VoiceStateUpdate) {
	return func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		log.Println("start TimeHandler")
		t.mu.Lock()
		defer t.mu.Unlock()

		user := v.UserID
		newCh := v.ChannelID

		if newCh != "" { // 채널 입장
			now := time.Now()
			chInfo, err := findChInfo(s, newCh)
			if err != nil {
				log.Println("채널 정보 가져오기 실패:", err)
				return
			}

			if _, ok := t.joinTimes[user]; !ok { // 음성채널 처음 입장
				t.joinTimes[user] = &joinInfo{
					joinedAt:  now,
					voiceChan: newCh,
				}

				err = t.sendMsg(s, newCh, fmt.Sprintf("<@%s>님이 %s에 %s로 오셨습니다.",
					user, now.Format(time.DateTime), chInfo.Name))
			} else { // 음성채널 이동
				beforeCh, err := findChInfo(s, t.joinTimes[user].voiceChan)
				if err != nil {
					log.Println("채널 정보 가져오기 실패:", err)
					return
				}

				t.joinTimes[user].voiceChan = newCh
				err = t.sendMsg(s, newCh, fmt.Sprintf("<@%s>님이 %s에 %s에서 %s로 이동하셨습니다.",
					user, now.Format(time.DateTime), beforeCh.Name, chInfo.Name))
			}
			if err != nil {
				log.Println("메시지 전송 실패:", err)
				return
			}
		}
		// 퇴장 newCh == ""
		if info, ok := t.joinTimes[user]; ok && newCh == "" {
			stayed := time.Since(info.joinedAt).Round(time.Second)
			chInfo, err := findChInfo(s, info.voiceChan)
			if err != nil {
				log.Println("채널 정보 가져오기 실패:", err)
				return
			}

			err = t.sendMsg(s, info.voiceChan, fmt.Sprintf("<@%s>님이 %s에서 퇴장하셨습니다. %s 동안 머물렀습니다.",
				user, chInfo.Name, stayed.String()))

			if err != nil {
				log.Println("메시지 전송 실패:", err)
				return
			}
			delete(t.joinTimes, user)
		}
	}
}

func findChInfo(s *discordgo.Session, chID string) (*discordgo.Channel, error) {
	channel, err := s.State.Channel(chID)
	if err != nil {
		channel, err = s.Channel(chID)
		if err != nil {
			return nil, fmt.Errorf("채널 정보 가져오기 실패: %w", err)
		}
	}
	return channel, nil
}

func (t *TimeHandler) sendMsg(s *discordgo.Session, voiceChID string, msg string) error {
	if ch, ok := t.VoiceToText[voiceChID]; ok {
		_, err := s.ChannelMessageSend(ch, msg)
		if err != nil {
			return err
		}
		return nil
	}

	// No Info
	return fmt.Errorf("채널 정보가 없습니다: %s", voiceChID)
}

func init() {
	m := map[string]string{
		"1198305978593062955": "1198305978140065848",
		"1198305978593062956": "1198305978140065849",
	}
	handler.Register(&TimeHandler{VoiceToText: m, joinTimes: map[string]*joinInfo{}})
}
