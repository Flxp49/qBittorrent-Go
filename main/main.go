package main

import (
	"log"
	"os"
	"time"

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

func main() {
	if godotenv.Load("../.env") != nil {
		log.Fatal("Error loading env file")
	}
	q, err := qbt.InitqBittorrent(os.Getenv("qbituser"), os.Getenv("qbitpass"), "http://localhost:3030", 3600)
	if err != nil {
		log.Fatal(err)
	}
	n := notion.InitNotionClient("Emad", os.Getenv("emadIntegSecret"), os.Getenv("emadDBID"))
	m := make(map[int]info) // map to maintain torrent state
	titles, err := n.QueryDB("Not started")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range titles.Results {
		_, exist := m[t.Properties.Tmdbid.Number]
		if !exist {
			sid, _ := q.InitSearch(t.Properties.Name.Title[0].Plain_text + " " + t.Properties.Year.Rich_text[0].Plain_text)
			time.Sleep(time.Second * 10)
			q.StopSearch(sid)
			res, _ := q.SearchJobResults(sid, 0)
			q.DeleteSearch(sid)
			q.AddTorrentDownload(res.Results[0].FileUrl, "D:\\Media\\Movies", "true", "true", "true")
			thash, _ := q.GetTorrentHash(res.Results[0].FileName, "all")
			m[t.Properties.Tmdbid.Number] = info{pageid: t.Pgid, title: t.Properties.Name.Title[0].Plain_text, year: t.Properties.Year.Rich_text[0].Plain_text, ttype: "movie", hash: thash, torrentName: res.Results[0].FileName}
			n.UpdateDownloadStatus(t.Pgid, "In progress", 0)
			time.Sleep(time.Second * 20)
			prog, _ := q.GetTorrentInfo(thash)
			log.Println(prog)
			n.UpdateDownloadStatus(t.Pgid, "In progress", float64(int(prog*100))/100)
		}
	}
}
