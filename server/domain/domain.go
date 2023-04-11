package domain

type PlatformOrgRepo struct {
	Platform string
	OrgRepo
}

type OrgRepo struct {
	Org  string
	Repo string
}

type Branch struct {
	Name string
	SHA  string
}

type RepoFileInfo struct {
	Name string
	Path string
	SHA  string
}

type RepoFile struct {
	RepoFileInfo

	// Allow empty file
	Content string
}
