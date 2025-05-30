package database

import (
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/database/supa"
	"github.com/supabase-community/supabase-go"
)

type Repository interface {
	// User 의 시간정보 가져오기
	// User 의 실제 공부시간 가져오기
}

func NewRepository() *supabase.Client {
	client := supa.Init()

	return client
}
