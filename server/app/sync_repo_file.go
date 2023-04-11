package app

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-repository-file/server/domain"
	"github.com/opensourceways/sync-repository-file/server/domain/repository"
)

type repoFileFilter struct {
	repo repository.RepoFile
}

func (s repoFileFilter) do(
	repo domain.PlatformOrgRepo,
	branch string,
	fileNames []string,
	allFiles []domain.RepoFileInfo,
) (r []string) {
	for _, fileName := range fileNames {
		cached, err := s.repo.FindFiles(repo, branch, fileName)
		if err != nil {
			logrus.Errorf("find cached files failed, err:%s", err)

			continue
		}

		if v := s.filterFile(fileName, cached, allFiles); len(v) > 0 {
			r = append(r, v...)
		}
	}

	return
}

func (s repoFileFilter) filterFile(fileName string, cached, files []domain.RepoFileInfo) []string {
	m := map[string]string{}
	for i := range cached {
		item := &cached[i]

		m[item.Path] = item.SHA
	}

	todo := make([]string, 0, len(files))
	for i := range files {
		item := &files[i]

		if item.Name != fileName {
			continue
		}

		if sha, ok := m[item.Path]; !ok || sha != item.SHA {
			todo = append(todo, item.Path)
		}
	}

	return todo
}
