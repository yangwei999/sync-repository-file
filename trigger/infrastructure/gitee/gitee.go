package gitee

import "github.com/opensourceways/robot-gitee-lib/client"

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

func (gp *giteePlatform) ListRepos(org string) ([]string, error) {
	repos, err := gp.cli.GetRepos(org)
	if err != nil || len(repos) == 0 {
		return nil, err
	}

	repoNames := make([]string, len(repos))

	for i := range repos {
		repoNames[i] = repos[i].Path
	}

	return repoNames, nil
}
