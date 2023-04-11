package domain

import "encoding/json"

// RepoFetchedEvent
type RepoFetchedEvent struct {
	Platform  string   `json:"platform"`
	Org       string   `json:"org"`
	Repo      string   `json:"repo"`
	FileNames []string `json:"file_names"`
}

func (t *RepoFetchedEvent) Message() ([]byte, error) {
	return json.Marshal(t)
}

func UnmarshalToRepoFetchedEvent(msg []byte) (r RepoFetchedEvent, err error) {
	err = json.Unmarshal(msg, &r)

	return
}
