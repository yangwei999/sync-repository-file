package messageimpl

type Config struct {
	Topics Topics `json:"topics"  required:"true"`
}

type Topics struct {
	RepoBranchFetched string `json:"repo_branch_fetched"  required:"true"`
	RepoFileFetched   string `json:"repo_file_fetched"    required:"true"`
}
