package codeplatform

type CodePlatform interface {
	Platform() string

	ListRepos(org string) ([]string, error)
}
