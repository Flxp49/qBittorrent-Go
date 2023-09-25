package main

import (
	"fmt"
	"os"

	"github.com/Flxp49/qbittorrent-notionwatchlist-sync/qbt"
	"github.com/joho/godotenv"
)

func main() {
	if godotenv.Load("../.env") != nil {
		fmt.Printf("Error loading env file")
		os.Exit(1)
	}
	_, err := qbt.InitqBittorrent(os.Getenv("qbituser"), os.Getenv("qbitpass"), "http://localhost:3030", 3600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
