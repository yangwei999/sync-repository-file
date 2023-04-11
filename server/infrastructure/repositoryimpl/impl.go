package repositoryimpl

import (
	"github.com/opensourceways/repo-file-cache/models"
	"github.com/opensourceways/repo-file-cache/sdk"

	"github.com/opensourceways/sync-repository-file/server/domain"
)

func NewRepoFileImpl(cfg *Config) *repoFileImpl {
	return &repoFileImpl{
		cli: sdk.NewSDK(cfg.Endpoint, 3),
	}
}

type repoFileImpl struct {
	cli *sdk.SDK
}

func (fc *repoFileImpl) SaveFile(repo domain.PlatformOrgRepo, branch string, file domain.RepoFile) error {
	opts := models.FileUpdateOption{
		Branch: models.Branch{
			Platform: repo.Platform,
			Org:      repo.Org,
			Repo:     repo.Repo,
			Branch:   branch,
		},
	}
	//opts.BranchSHA = branch.SHA

	opts.Files = []models.File{
		{
			Path:    models.FilePath(file.Path),
			SHA:     file.SHA,
			Content: file.Content,
		},
	}

	return fc.cli.SaveFiles(opts)
}

func (fc *repoFileImpl) FindFiles(repo domain.PlatformOrgRepo, branch, fileName string) (
	[]domain.RepoFileInfo, error,
) {
	v, err := fc.cli.GetFiles(
		models.Branch{
			Platform: repo.Platform,
			Org:      repo.Org,
			Repo:     repo.Repo,
			Branch:   branch,
		},
		fileName,
		true,
	)
	if err != nil {
		return nil, err
	}

	r := make([]domain.RepoFileInfo, len(v.Files))

	for i, item := range v.Files {
		r[i].Path = item.Path.FullPath()
		r[i].SHA = item.SHA
	}

	return r, nil
}
