package main

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

func main() {
	supabaseUrl := "https://rojqpylyzaaryefgalrp.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJvanFweWx5emFhcnllZmdhbHJwIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDIxODcwNjYsImV4cCI6MjA1Nzc2MzA2Nn0.4CbjARN73KTJVma3Yarf2CNT6FyPw2qGW9ENyv_f5Ns"
	// supabaseKey := "YOUR_SUPABASE_ANON_KEY" // or YOUR_SUPABASE_SERVICE_ROLE_KEY for admin privileges
	supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
	if err != nil {
		panic(err)
	}

	var results []map[string]interface{}
	data, count, err := supabaseClient.From("notes").Select("*", "", false).Execute()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &results)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("Count:", count)
	fmt.Println(results)

}
