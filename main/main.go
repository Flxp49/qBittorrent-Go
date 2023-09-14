package main

import (
	"fmt"
	"os"

	"github.com/Flxp49/qbittorrent-notionwatchlist-sync/qbt"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Printf("Hello")
	if godotenv.Load("../.env") != nil {
		fmt.Printf("Error loading env file")
		os.Exit(1)
	}
	qbt.InitqBittorrent(os.Getenv("user"), os.Getenv("pass"), "http://localhost:3030")
}
