package messageimpl

import (
	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/sync-repository-file/server/domain/message"
)

const (
	HeaderKey                    = "headerKey"
	HeaderValueRepoFetched       = "handleRepoFetched"
	HeaderValueRepoBranchFetched = "handleRepoBranchFetched"
	HeaderValueRepoFileFetched   = "handleRepoFileFetched"
)

func NewRepoFileMessage(cfg *Config) *repoFileMessage {
	return &repoFileMessage{
		topics: cfg.Topics,
	}
}

type repoFileMessage struct {
	topics Topics
}

func (p *repoFileMessage) SendRepoBranchFetchedEvent(e message.Message) error {
	header := map[string]string{
		HeaderKey: HeaderValueRepoBranchFetched,
	}

	return send(p.topics.RepoBranchFetched, header, e)
}

func (p *repoFileMessage) SendRepoFileFetchedEvent(e message.Message) error {
	header := map[string]string{
		HeaderKey: HeaderValueRepoFileFetched,
	}

	return send(p.topics.RepoFileFetched, header, e)
}

func send(topic string, header map[string]string, v message.Message) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, header, body)
}
