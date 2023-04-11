package app

import (
	"github.com/opensourceways/sync-repository-file/trigger/domain"
	"github.com/opensourceways/sync-repository-file/trigger/domain/codeplatform"
	"github.com/opensourceways/sync-repository-file/trigger/domain/message"
)

type RepoService interface {
	FetchRepo(codeplatform.CodePlatform, *CmdToFetchRepo) error
}

func NewRepoService(message message.Repo) repoService {
	return repoService{
		message: message,
	}
}

type repoService struct {
	message message.Repo
}

func (s repoService) FetchRepo(p codeplatform.CodePlatform, cmd *CmdToFetchRepo) error {
	v, err := p.ListRepos(cmd.Org)
	if err != nil {
		return err
	}

	v = cmd.filterRepos(v)

	task := domain.RepoFetchedEvent{
		Platform:  p.Platform(),
		Org:       cmd.Org,
		FileNames: cmd.FileNames,
	}

	for i := range v {
		task.Repo = v[i]

		if err := s.message.SendRepoFetchedEvent(&task); err != nil {
			return err
		}
	}

	return nil
}
