package messageimpl

import (
	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/sync-repository-file/server/domain/message"
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
	return send(p.topics.RepoBranchFetched, e)
}

func (p *repoFileMessage) SendRepoFileFetchedEvent(e message.Message) error {
	return send(p.topics.RepoFileFetched, e)
}

func send(topic string, v message.Message) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, nil, body)
}
