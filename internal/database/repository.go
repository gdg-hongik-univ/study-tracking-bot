package database

import (
	"github.com/gdg-hongik-univ/study-tracking-bot/internal/database/supa"
	"github.com/supabase-community/supabase-go"
)

type Repository interface {
}

func NewRepository() *supabase.Client {
	client := supa.Init()

	return client
}
