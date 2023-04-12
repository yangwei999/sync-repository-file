package domain

import "encoding/json"

// RepoBranchFetchedEvent
type RepoBranchFetchedEvent struct {
	Platform  string   `json:"platform"`
	Org       string   `json:"org"`
	Repo      string   `json:"repo"`
	Branch    string   `json:"branch"`
	BranchSHA string   `json:"branch_sha"`
	FileNames []string `json:"file_names"`
}

func (t *RepoBranchFetchedEvent) Message() ([]byte, error) {
	return json.Marshal(t)
}

func UnmarshalToRepoBranchFetchedEvent(msg []byte) (r RepoBranchFetchedEvent, err error) {
	err = json.Unmarshal(msg, &r)

	return
}

// RepoFileFetchedEvent
type RepoFileFetchedEvent struct {
	Platform  string `json:"platform"`
	Org       string `json:"org"`
	Repo      string `json:"repo"`
	Branch    string `json:"branch"`
	BranchSHA string `json:"branch_sha"`
	FilePath  string `json:"path"`
}

func (t *RepoFileFetchedEvent) Message() ([]byte, error) {
	return json.Marshal(t)
}

func UnmarshalToRepoFileFetchedEvent(msg []byte) (r RepoFileFetchedEvent, err error) {
	err = json.Unmarshal(msg, &r)

	return
}
