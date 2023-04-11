package message

type Message interface {
	Message() ([]byte, error)
}

type Repo interface {
	SendRepoFetchedEvent(Message) error
}
