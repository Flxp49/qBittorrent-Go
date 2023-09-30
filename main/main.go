package main

import (
	"log"
	"os"
	"sync"

	"github.com/Flxp49/qbittorrent-notionwatchlist-sync/notion"
	"github.com/Flxp49/qbittorrent-notionwatchlist-sync/qbt"
	"github.com/joho/godotenv"
)

// active torrent struct containing info mapped to tmdb id as key
type info struct {
	pageid      string
	hash        string
	title       string
	year        string
	ttype       string
	torrentName string
}

type mainInstance struct {
	cm      map[int]info    // map to maintain current torrent
	maplock *sync.RWMutex   //lock for map
	qbt     *qbt.QBittorent // qbt instance
	nt      *notion.NotionClient
}

func main() {
	if godotenv.Load("../.env") != nil {
		log.Fatal("Error loading env file")
	}
	var lock sync.Mutex
	qbit, err := qbt.InitqBittorrent(os.Getenv("qbituser"), os.Getenv("qbitpass"), "http://localhost:3030", 3600, &lock)
	if err != nil {
		log.Fatal("Failed to authorize qbt", err)
	}
	var maplock sync.RWMutex
	m := make(map[int]info)
	mi := mainInstance{cm: m, maplock: &maplock, qbt: qbit, nt: notion.InitNotionClient("Emad", os.Getenv("emadIntegSecret"), os.Getenv("emadDBID"))}
	Init(mi)
}
