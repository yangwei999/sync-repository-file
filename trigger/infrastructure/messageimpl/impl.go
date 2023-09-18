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
	header := map[string]string{
		"header_key": "handleRepoFetched",
	}

	return send(p.topics.RepoFetched, header, e)
}

func send(topic string, header map[string]string, v message.Message) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kafka.Publish(topic, header, body)
}
