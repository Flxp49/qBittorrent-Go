package notion

import (
	"bytes"
	"encoding/json"
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

func parseJson(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

func (n *notionClient) performReq(method string, endpoint string, data []byte) (*http.Response, []byte, error) {
	n.req.Method = method
	n.req.URL, _ = url.Parse("https://api.notion.com" + "/" + endpoint)
	if method == "POST" || method == "PATCH" {
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

// QueryDB Response struct
type QueryDB struct {
	Results []struct {
		Pgid       string `json:"id"`
		Properties struct {
			Tmdbid struct {
				Number int `json:"number"`
			} `json:"ID"`
			Name struct {
				Title []struct {
					Plain_text string `json:"plain_text"`
				} `json:"title"`
			}
			Year struct {
				Rich_text []struct {
					Plain_text string `json:"plain_text"`
				} `json:"rich_text"`
			}
		} `json:"properties"`
	} `json:"results"`
}

// QueryDB payload struct
type QueryDBPayload struct {
	Filter struct {
		And []struct {
			Property string `json:"property"`
			Checkbox *struct {
				Equals bool `json:"equals"`
			} `json:"checkbox,omitempty"`
			Select *struct {
				Equals string `json:"equals"`
			} `json:"select,omitempty"`
		} `json:"and"`
	} `json:"filter"`
}

// This funtion returns the list of tmdb ids and pageids that are checked as downlaod
//
// status - fetch records with download status as "Not started", "In progress" or "Done"
func (n *notionClient) QueryDB(status string) (QueryDB, error) {
	payload := QueryDBPayload{
		Filter: struct {
			And []struct {
				Property string `json:"property"`
				Checkbox *struct {
					Equals bool `json:"equals"`
				} `json:"checkbox,omitempty"`
				Select *struct {
					Equals string `json:"equals"`
				} `json:"select,omitempty"`
			} `json:"and"`
		}{
			And: []struct {
				Property string `json:"property"`
				Checkbox *struct {
					Equals bool `json:"equals"`
				} `json:"checkbox,omitempty"`
				Select *struct {
					Equals string `json:"equals"`
				} `json:"select,omitempty"`
			}{
				// {
				// 	Property: "Download Status",
				// 	Select: &struct {
				// 		Equals string `json:"equals"`
				// 	}{
				// 		Equals: status,
				// 	},
				// },
				{
					Property: "Type",
					Select: &struct {
						Equals string `json:"equals"`
					}{
						Equals: "Movie",
					},
				},
				{
					Property: "Download",
					Checkbox: &struct {
						Equals bool `json:"equals"`
					}{
						Equals: true,
					},
				},
			},
		},
	}
	data, _ := json.Marshal(payload)
	_, body, err := n.performReq("POST", fmt.Sprintf("v1/databases/%s/query", n.dbid), data)
	if err != nil {
		return QueryDB{}, err
	}
	var qDb QueryDB
	err = parseJson(body, &qDb)
	if err != nil {
		return QueryDB{}, err
	}
	return qDb, nil
}

type statusMap struct {
	name  string
	color string
}

var sMap = map[string]statusMap{
	"error":              {name: "🔴 Error", color: "red"},
	"missingFiles":       {name: "🔴 Missing files", color: "red"},
	"uploading":          {name: "🟢 Downloaded", color: "green"},
	"pausedUP":           {name: "🟢 Downloaded", color: "green"},
	"queuedUP":           {name: "🟢 Downloaded", color: "green"},
	"stalledUP":          {name: "🟢 Downloaded", color: "green"},
	"checkingUP":         {name: "🟢 Downloaded", color: "green"},
	"forcedUP":           {name: "🟢 Downloaded", color: "green"},
	"allocating":         {name: "🟡 Allocating", color: "yellow"},
	"downloading":        {name: "🔵 Downloading", color: "blue"},
	"metaDL":             {name: "🟡 Downloading metadata", color: "yellow"},
	"pausedDL":           {name: "🟠 Paused", color: "orange"},
	"queuedDL":           {name: "🟡 Queued", color: "yellow"},
	"stalledDL":          {name: "🟠 Stalled", color: "orange"},
	"checkingDL":         {name: "🟡 Checking download", color: "yellow"},
	"forcedDL":           {name: "🔵 Downloading", color: "blue"},
	"checkingResumeData": {name: "🟡 Checking resume data", color: "yellow"},
	"moving":             {name: "🟡 Moving ", color: "yellow"},
	"unknown":            {name: "⭕ Unknown ", color: "gray"},
}

// Payload struct for UpdateDownloadStatus
type UpdateDownloadStatus struct {
	Properties struct {
		DStatus struct {
			Select struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"select"`
		} `json:"Download Status"`
		Dprogress struct {
			Number float64 `json:"number"`
		} `json:"Download Progress"`
	} `json:"properties"`
}

// this UpdateDownloadStatus function updates the "Download Status" prop
//
// id - page id to update
//
// status - "Not started" or "In progress" or "Done"
//
// val - Download progress [0,1]
func (n *notionClient) UpdateDownloadStatus(id string, status string, val float64) error {
	UpdateDownloadStatus := UpdateDownloadStatus{Properties: struct {
		DStatus struct {
			Select struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"select"`
		} `json:"Download Status"`
		Dprogress struct {
			Number float64 `json:"number"`
		} `json:"Download Progress"`
	}{
		DStatus: struct {
			Select struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"select"`
		}{
			Select: struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			}{
				Name:  sMap[status].name,
				Color: sMap[status].color,
			},
		},
		Dprogress: struct {
			Number float64 `json:"number"`
		}{
			Number: val,
		},
	},
	}

	data, err := json.Marshal(UpdateDownloadStatus)
	if err != nil {
		return err
	}
	_, _, err = n.performReq("PATCH", fmt.Sprintf("v1/pages/%s", id), data)
	if err != nil {
		return err
	}
	return nil
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