package main

import (
	"github.com/opensourceways/sync-repository-file/server/app"
	"github.com/opensourceways/sync-repository-file/server/domain"
)

// msgOfRepoFetched
type msgOfRepoFetched struct {
	Platform  string   `json:"platform"`
	Org       string   `json:"org"`
	Repo      string   `json:"repo"`
	FileNames []string `json:"file_names"`
}

func (msg *msgOfRepoFetched) toCmd() app.CmdToFetchRepoBranch {
	return app.CmdToFetchRepoBranch{
		OrgRepo: domain.OrgRepo{
			Org:  msg.Org,
			Repo: msg.Repo,
		},
		FileNames: msg.FileNames,
	}
}

// cmdToFetchRepoFile
func cmdToFetchRepoFile(data []byte) (
	cmd app.CmdToFetchRepoFile, platform string, err error,
) {
	v, err := domain.UnmarshalToRepoBranchFetchedEvent(data)
	if err != nil {
		return
	}

	cmd.Org = v.Org
	cmd.Repo = v.Repo
	cmd.Branch = v.Branch
	cmd.FileNames = v.FileNames

	platform = v.Platform

	return
}

// cmdToFetchFileContent
func cmdToFetchFileContent(data []byte) (
	cmd app.CmdToFetchFileContent, platform string, err error,
) {
	v, err := domain.UnmarshalToRepoFileFetchedEvent(data)
	if err != nil {
		return
	}

	cmd.Org = v.Org
	cmd.Repo = v.Repo
	cmd.Branch = v.Branch
	cmd.FilePath = v.FilePath

	platform = v.Platform

	return
}
