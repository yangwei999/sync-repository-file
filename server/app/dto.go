package app

import "github.com/opensourceways/sync-repository-file/server/domain"

type CmdToFetchRepoBranch struct {
	domain.OrgRepo

	FileNames []string
}

type CmdToFetchRepoFile struct {
	domain.OrgRepo

	Branch    string
	FileNames []string
}

type CmdToFetchFileContent struct {
	domain.OrgRepo

	Branch   string
	FilePath string
}
