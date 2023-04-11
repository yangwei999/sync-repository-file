package codeplatform

import "github.com/opensourceways/sync-repository-file/server/domain"

type CodePlatform interface {
	Platform() string

	ListBranches(domain.OrgRepo) ([]domain.Branch, error)

	ListFiles(repo domain.OrgRepo, branch string) ([]domain.RepoFileInfo, error)

	GetFile(repo domain.OrgRepo, branch, path string) (domain.RepoFile, error)
}
