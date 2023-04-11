package messageimpl

import (
	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/sync-repository-file/trigger/domain/message"
)

func NewRepoMessage(cfg *Config) *repoMessage {
	return &repoMessage{
		topics: cfg.Topics,
	}
}

type repoMessage struct {
	topics Topics
}

func (p *repoMessage) SendRepoFetchedEvent(e message.Message) error {
	return send(p.topics.RepoFetched, e)
}

func send(topic string, v message.Message) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, nil, body)
}
