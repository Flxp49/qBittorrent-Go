package main

import (
	"log"
	"time"
)

func Init(mi mainInstance) {
	titles, err := mi.nt.QueryDB("Not started")
	if err != nil {
		log.Fatal(err)
	}
	// fetchTitles(&mi)
	titles, err = mi.nt.QueryDB("Not started") //fetch titles from the db which have "Not started" status
	if err != nil {
		log.Panicln(err)
	}
	mi.maplock.Lock()
	for _, t := range titles.Results {
		_, exist := mi.cm[t.Properties.Tmdbid.Number]
		if !exist {
			sid, _ := mi.qbt.InitSearch(t.Properties.Name.Title[0].Plain_text + " " + t.Properties.Year.Rich_text[0].Plain_text)
			time.Sleep(time.Second * 10)
			mi.qbt.StopSearch(sid)
			res, _ := mi.qbt.SearchJobResults(sid, 0)
			mi.qbt.DeleteSearch(sid)
			mi.qbt.AddTorrentDownload(res.Results[0].FileUrl, "D:\\Media\\Movies", "true", "true", "true")
			thash, _ := mi.qbt.GetTorrentHash(res.Results[0].FileName, "all")
			mi.cm[t.Properties.Tmdbid.Number] = info{pageid: t.Pgid, title: t.Properties.Name.Title[0].Plain_text, year: t.Properties.Year.Rich_text[0].Plain_text, ttype: "movie", hash: thash, torrentName: res.Results[0].FileName}
			mi.maplock.Unlock()
			prog, _ := mi.qbt.GetTorrentInfo(thash)
			mi.nt.UpdateDownloadStatus(t.Pgid, prog[0].State, float64(int(prog[0].Progress*100))/100)
			time.Sleep(time.Second * 15)
			prog, _ = mi.qbt.GetTorrentInfo(thash)
			mi.nt.UpdateDownloadStatus(t.Pgid, prog[0].State, float64(int(prog[0].Progress*100))/100)
		}
	}
	mi.maplock.Unlock()
}

// this function fetches the titles from watchlist
// func fetchTitles(mi *mainInstance) error {

// 	if err != nil {
// 		log.Println("Error querying notion DB", err)
// 		return err
// 	}
// 	mi.maplock.RLock()
// 	mi.maplock.Lock()
// 	defer mi.maplock.RUnlock()

// }
