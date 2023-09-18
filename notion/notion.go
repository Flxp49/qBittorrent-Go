package notion

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type notionClient struct {
	user   string
	secret string
	dbid   string
	req    *http.Request
}

func (n *notionClient) performReq(method string, endpoint string, data []byte) (*http.Response, []byte, error) {
	n.req.Method = method
	n.req.URL, _ = url.Parse("https://api.notion.com" + "/" + endpoint)
	if method == "POST" {
		n.req.Body = io.NopCloser(bytes.NewBuffer(data))
		n.req.ContentLength = int64(len(data))
	} else {
		n.req.Body = nil
		n.req.ContentLength = 0
	}
	resp, err := http.DefaultClient.Do(n.req)
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

// This funtion returns the list of tmdb ids that are under the watching section and checked as downlaod
func (n *notionClient) queryDB() {

}

// constructor
func InitNotionClient(username string, secret string, dbid string) *notionClient {
	n := &notionClient{user: username, secret: secret, dbid: dbid}
	n.req, _ = http.NewRequest("", "", nil)
	n.req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.secret))
	n.req.Header.Add("Notion-Version", "2022-06-28")
	n.req.Header.Add("Content-Type", "application/json")
	return n
}
