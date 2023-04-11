package main

import (
	"github.com/opensourceways/server-common-lib/utils"

	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/sync-repository-file/server/infrastructure/gitee"
	"github.com/opensourceways/sync-repository-file/server/infrastructure/messageimpl"
	"github.com/opensourceways/sync-repository-file/server/infrastructure/repositoryimpl"
)

func loadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.setDefault()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

// Config
type Config struct {
	Kafka      kafka.Config          `json:"kafka"                required:"true"`
	Gitee      gitee.Config          `json:"gitee"                required:"true"`
	Topics     Topics                `json:"topics_to_subscribe"  required:"true"`
	GroupName  string                `json:"group_name"           required:"true"`
	Message    messageimpl.Config    `json:"message"              required:"true"`
	Repository repositoryimpl.Config `json:"repository"           required:"true"`
}

type Topics struct {
	RepoFetched       string `json:"repo_fetched"         required:"true"`
	RepoBranchFetched string `json:"repo_branch_fetched"  required:"true"`
	RepoFileFetched   string `json:"repo_file_fetched"    required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Kafka,
		&cfg.Gitee,
		&cfg.Topics,
		&cfg.Message,
		&cfg.Repository,
	}
}

func (cfg *Config) setDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
