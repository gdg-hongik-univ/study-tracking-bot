package supa

import (
	"github.com/supabase-community/supabase-go"
	"os"
)

func Init() *supabase.Client {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		panic(err)
	}
	return client
}
