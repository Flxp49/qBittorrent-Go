package qbt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func parseJson(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

type qBittorent struct {
	username       string
	password       string
	host           string
	req            *http.Request
	cookieTime     time.Time
	sessionTimeout time.Duration
}

func (q *qBittorent) performReq(method string, endpoint string, data []byte) (*http.Response, []byte, error) {
	if time.Since(q.cookieTime) >= q.sessionTimeout {
		err := q.requestAuth()
		if err != nil {
			return nil, nil, err
		}
	}
	q.req.Method = method
	q.req.URL, _ = url.Parse(q.host + "/" + endpoint)
	if method == "POST" {
		q.req.Body = io.NopCloser(bytes.NewBuffer(data))
		q.req.ContentLength = int64(len(data))
	} else {
		q.req.Body = nil
		q.req.ContentLength = 0
	}
	resp, err := http.DefaultClient.Do(q.req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, body, nil
}

func (q *qBittorent) requestAuth() error {

	data := []byte(fmt.Sprintf("username=%s&password=%s", q.username, q.password))
	req, _ := http.NewRequest("POST", q.host+"/"+"api/v2/auth/login", bytes.NewBuffer(data))
	req.Header.Add("Referer", q.host)
	req.Header.Add("Host", q.host)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if len(resp.Cookies()) != 0 {
		q.req.AddCookie(resp.Cookies()[0])
		q.cookieTime = time.Now()
		return nil
	}
	return err
}

// initSearch: start search job result struct
type initSearch struct {
	Id int `json:"id"`
}

// Starts torrent search and returns id of the search job
//
// pattern - string to search torrents by
func (q *qBittorent) InitSearch(pattern string) (int, error) {
	data := []byte(fmt.Sprintf("pattern=%s&plugins=enabled&category=all", pattern))
	_, body, err := q.performReq("POST", "api/v2/search/start", data)
	if err != nil {
		return -1, err
	}
	var sji initSearch
	err = parseJson(body, &sji)
	if err != nil {
		return -1, err
	}
	return sji.Id, nil
}

// Stops torrent search job
//
// id - search id of search job to stop
func (q *qBittorent) StopSearch(id int) error {
	data := []byte(fmt.Sprintf("id=%d", id))
	_, _, err := q.performReq("POST", "api/v2/search/stop", data)
	return err
}

// Deletes torrrent search
//
// id - search id of search job to delete
func (q *qBittorent) DeleteSearch(id int) error {
	data := []byte(fmt.Sprintf("id=%d", id))
	_, _, err := q.performReq("POST", "api/v2/search/delete", data)
	return err
}

// searchJobResults: torrent search job result struct
type searchJobResults struct {
	Results []struct {
		DescrLink  string `json:"descrLink"`
		FileName   string `json:"fileName"`
		FileSize   int    `json:"fileSize"`
		FileUrl    string `json:"fileUrl"`
		NbLeechers int    `json:"nbLeechers"`
		NbSeeders  int    `json:"nbSeeders"`
		SiteUrl    string `json:"siteUrl"`
	} `json:"results"`
	Status string `json:"status"`
	Total  int    `json:"total"`
}

// Get torrrent search job results
//
// id - search id
//
// limit - search results limit, 0 => no limit
func (q *qBittorent) SearchJobResults(id int, limit int) (searchJobResults, error) {
	data := []byte(fmt.Sprintf("id=%d", id))
	_, body, err := q.performReq("POST", "api/v2/search/results", data)
	if err != nil {
		return searchJobResults{}, err
	}
	var sjr searchJobResults
	err = parseJson(body, &sjr)
	if err != nil {
		return searchJobResults{}, err
	}
	return sjr, nil
}

// Add torrrent to download
//
// urls - URLs separated with newlines
//
// savepath - Download folder
//
// root_folder - Create the root folder. Possible values are true, false
//
// sequentialDownload - Enable sequential download. Possible values are true
//
// firstLastPiecePrio - Prioritize download first last piece. Possible values are true
func (q *qBittorent) AddTorrentDownload(urls string, savepath string, root_folder string, sequentialDownload string, firstLastPiecePrio string) error {
	data := []byte(fmt.Sprintf("urls=%s&savepath=%s&root_folder=%s&sequentialDownload=%s&firstLastPiecePrio=%s", urls, savepath, root_folder, sequentialDownload, firstLastPiecePrio))
	_, _, err := q.performReq("POST", "api/v2/torrents/add", data)
	return err
}

// getTorrentHashResult: torrent hash result struct
type getTorrentHashResult []struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
}

//	Get torrent hash by name, returns torrent hash
//
// name - Name of the torrent to fetch hash of
//
// filter - Filter torrent list by state. Allowed state filters: all, downloading, seeding, completed, paused, active, inactive, resumed, stalled, stalled_uploading, stalled_downloading, errored
func (q *qBittorent) GetTorrentHash(name string, filter string) (string, error) {
	url := fmt.Sprintf("api/v2/torrents/info?filter=%s", filter)
	_, body, err := q.performReq("GET", url, nil)
	if err != nil {
		return "", err
	}
	var gthr getTorrentHashResult
	err = parseJson(body, &gthr)
	if err != nil {
		return "", err
	}

	for _, t := range gthr {
		if t.Name == name {
			return t.Hash, nil
		}
	}
	return "", nil
}

// GetTorrentInfoResult: torrent info result struct
type GetTorrentInfoResult []struct {
	Progress float64 `json:"progress"`
	State    string  `json:"state"`
}

//	Get torrent info by hash, returns torrent progress
//
// name - Name of the torrent to fetch hash of
//
// filter - Filter torrent list by state. Allowed state filters: all, downloading, seeding, completed, paused, active, inactive, resumed, stalled, stalled_uploading, stalled_downloading, errored
func (q *qBittorent) GetTorrentInfo(hash string) (GetTorrentInfoResult, error) {
	url := fmt.Sprintf("api/v2/torrents/info?hashes=%s", hash)
	_, body, err := q.performReq("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var gtir GetTorrentInfoResult
	err = parseJson(body, &gtir)
	if err != nil {
		return nil, err
	}
	return gtir, nil
}

// Delete torrrent by hash
//
// hashes - Hash of torrent to delete
//
// deleteFiles - If set to true, the downloaded data will also be deleted, otherwise has no effect
func (q *qBittorent) DeleteTorrent(hashes string, deletefiles string) error {
	data := []byte(fmt.Sprintf("hashes=%s", hashes))
	_, _, err := q.performReq("POST", "api/v2/torrents/delete", data)
	return err
}

// constructor
func InitqBittorrent(username string, password string, host string, timeout int) (*qBittorent, error) {
	q := &qBittorent{username: username, password: password, host: host, sessionTimeout: time.Second * time.Duration(timeout)}
	q.req, _ = http.NewRequest("", "", nil)
	err := q.requestAuth()
	if err != nil {
		return nil, err
	}
	q.req.Header.Add("Referer", q.host)
	q.req.Header.Add("Host", q.host)
	q.req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return q, nil
}
