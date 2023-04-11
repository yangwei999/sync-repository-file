package main

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-repository-file/trigger/app"
	"github.com/opensourceways/sync-repository-file/trigger/domain/codeplatform"
)

type triggerImpl struct {
	platforms map[string]codeplatform.CodePlatform
	service   app.RepoService
}

func (impl *triggerImpl) kickoff(v repos) {
	for i := range v {
		impl.start(&v[i])
	}
}

func (impl *triggerImpl) start(cfg *repoConfig) {
	p, ok := impl.platforms[cfg.Platform]
	if !ok {
		logrus.Errorf("unknown platform:%s", cfg.Platform)

		return
	}

	cmd := app.CmdToFetchRepo{
		FileNames: cfg.FileNames,
	}

	for i := range cfg.OrgRepos {
		item := &cfg.OrgRepos[i]

		cmd.Org = item.Org
		cmd.Repos = item.Repos
		cmd.ExcludedRepos = item.ExcludedRepos

		if err := impl.service.FetchRepo(p, &cmd); err != nil {
			logrus.Errorf("fetch repo failed for org: %s", item.Org)
		}
	}

}
