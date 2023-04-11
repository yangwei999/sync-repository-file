package app

import "k8s.io/apimachinery/pkg/util/sets"

type CmdToFetchRepo struct {
	Org           string
	Repos         []string
	ExcludedRepos []string
	FileNames     []string
}

func (cmd *CmdToFetchRepo) filterRepos(repos []string) []string {
	if len(cmd.Repos) > 0 {
		return sets.NewString(repos...).Intersection(
			sets.NewString(cmd.Repos...),
		).UnsortedList()
	}

	if len(cmd.ExcludedRepos) > 0 {
		return sets.NewString(repos...).Difference(
			sets.NewString(cmd.ExcludedRepos...),
		).UnsortedList()
	}

	return repos
}
