package gitee

import (
	"fmt"
	"path"

	"github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/sync-repository-file/server/domain"
	"github.com/opensourceways/sync-repository-file/utils"
)

func NewGiteePlatform(cfg *Config) *giteePlatform {
	return &giteePlatform{
		cli: client.NewClient(
			func() []byte {
				return []byte(cfg.Token)
			},
		),
	}
}

type giteePlatform struct {
	cli client.Client
}

func (gp *giteePlatform) Platform() string {
	return "gitee"
}

func (gp *giteePlatform) ListBranches(repo domain.OrgRepo) ([]domain.Branch, error) {
	branches, err := gp.cli.GetRepoAllBranch(repo.Org, repo.Repo)
	if err != nil || len(branches) == 0 {
		return nil, fmt.Errorf("%s of %s in %s", err.Error(), repo.Repo, repo.Org)
	}

	infos := make([]domain.Branch, len(branches))

	for i := range branches {
		item := &branches[i]

		infos[i].Name = item.GetName()
		infos[i].SHA = item.GetCommit().GetSha()
	}

	return infos, err
}

func (gp *giteePlatform) ListFiles(repo domain.OrgRepo, branch string) (
	[]domain.RepoFileInfo, error,
) {
	var trees gitee.Tree
	var err error

	err = utils.Retry(func() error {
		trees, err = gp.cli.GetDirectoryTree(repo.Org, repo.Repo, branch, 1)

		return err
	})

	if err != nil || len(trees.Tree) == 0 {
		return nil, err
	}

	files := make([]domain.RepoFileInfo, len(trees.Tree))

	for i := range trees.Tree {
		item := &trees.Tree[i]

		files[i].Name = path.Base(item.Path)
		files[i].Path = item.Path
		files[i].SHA = item.Sha
	}

	return files, nil
}

func (gp *giteePlatform) GetFile(repo domain.OrgRepo, branch, filepath string) (
	r domain.RepoFile, err error,
) {
	content, err := gp.cli.GetPathContent(repo.Org, repo.Repo, filepath, branch)
	if err == nil {
		r.Name = path.Base(filepath)
		r.Path = filepath
		r.SHA = content.Sha
		r.Content = content.Content
	}

	return
}
