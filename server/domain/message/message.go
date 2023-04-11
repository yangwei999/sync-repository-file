package message

type Message interface {
	Message() ([]byte, error)
}

type RepoFile interface {
	SendRepoBranchFetchedEvent(Message) error
	SendRepoFileFetchedEvent(Message) error
}
