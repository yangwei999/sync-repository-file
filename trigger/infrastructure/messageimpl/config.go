package messageimpl

type Config struct {
	Topics Topics `json:"topics"  required:"true"`
}

type Topics struct {
	RepoFetched string `json:"repo_fetched"  required:"true"`
}
