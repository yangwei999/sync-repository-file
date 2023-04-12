package app

import (
	"github.com/opensourceways/sync-repository-file/server/domain"
	"github.com/opensourceways/sync-repository-file/server/domain/codeplatform"
	"github.com/opensourceways/sync-repository-file/server/domain/message"
	"github.com/opensourceways/sync-repository-file/server/domain/repository"
)

type RepoFileService interface {
	FetchRepoBranch(codeplatform.CodePlatform, *CmdToFetchRepoBranch) error
	FetchRepoFile(codeplatform.CodePlatform, *CmdToFetchRepoFile) error
	FetchFileContent(codeplatform.CodePlatform, *CmdToFetchFileContent) error
}

func NewRepoFileService(
	repo repository.RepoFile,
	message message.RepoFile,
) repoFileService {
	return repoFileService{
		repo:           repo,
		message:        message,
		repoFileFilter: repoFileFilter{repo},
	}
}

type repoFileService struct {
	repo           repository.RepoFile
	message        message.RepoFile
	repoFileFilter repoFileFilter
}

func (s repoFileService) FetchRepoBranch(
	p codeplatform.CodePlatform,
	cmd *CmdToFetchRepoBranch,
) error {
	v, err := p.ListBranches(cmd.OrgRepo)
	if err != nil {
		return err
	}

	task := domain.RepoBranchFetchedEvent{
		Platform:  p.Platform(),
		Org:       cmd.Org,
		Repo:      cmd.Repo,
		FileNames: cmd.FileNames,
	}

	for i := range v {
		task.Branch = v[i].Name
		task.BranchSHA = v[i].SHA

		if err := s.message.SendRepoBranchFetchedEvent(&task); err != nil {
			return err
		}
	}

	return nil
}

func (s repoFileService) FetchRepoFile(
	p codeplatform.CodePlatform,
	cmd *CmdToFetchRepoFile,
) error {
	v, err := p.ListFiles(cmd.OrgRepo, cmd.Branch.Name)
	if err != nil {
		return err
	}

	files := s.repoFileFilter.do(
		domain.PlatformOrgRepo{
			Platform: p.Platform(),
			OrgRepo: domain.OrgRepo{
				Org:  cmd.Org,
				Repo: cmd.Repo,
			},
		},
		cmd.Branch.Name, cmd.FileNames, v,
	)
	if len(files) == 0 {
		return nil
	}

	task := domain.RepoFileFetchedEvent{
		Platform:  p.Platform(),
		Org:       cmd.Org,
		Repo:      cmd.Repo,
		Branch:    cmd.Branch.Name,
		BranchSHA: cmd.Branch.SHA,
	}

	for _, path := range files {
		task.FilePath = path

		if err := s.message.SendRepoFileFetchedEvent(&task); err != nil {
			return err
		}
	}

	return nil
}

func (s repoFileService) FetchFileContent(
	p codeplatform.CodePlatform,
	cmd *CmdToFetchFileContent,
) error {
	v, err := p.GetFile(cmd.OrgRepo, cmd.Branch.Name, cmd.FilePath)
	if err != nil {
		return err
	}

	return s.repo.SaveFile(
		domain.PlatformOrgRepo{
			Platform: p.Platform(),
			OrgRepo:  cmd.OrgRepo,
		},
		cmd.Branch,
		v,
	)
}
