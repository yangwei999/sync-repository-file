package main

import (
	"context"
	"encoding/json"
	"errors"

	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/sync-repository-file/server/app"
	"github.com/opensourceways/sync-repository-file/server/domain/codeplatform"
)

const (
	HeaderKey                  = "header_key"
	HeaderKeyRepoFetched       = "handleRepoFetched"
	HeaderKeyRepoBranchFetched = "handleRepoBranchFetched"
	HeaderKeyRepoFileFetched   = "handleRepoFileFetched"
)

type server struct {
	service   app.RepoFileService
	platforms map[string]codeplatform.CodePlatform
}

func (s *server) run(ctx context.Context, cfg *Config) error {
	if err := s.subscribe(cfg); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *server) subscribe(cfg *Config) error {
	return kafka.Subscribe(cfg.GroupName, s.handleByHeader, []string{
		cfg.Topics.RepoFetched,
		cfg.Topics.RepoBranchFetched,
		cfg.Topics.RepoFileFetched,
	})
}

func (s *server) handleByHeader(data []byte, header map[string]string) error {
	switch header[HeaderKey] {
	case HeaderKeyRepoFetched:
		return s.handleRepoFetched(data, header)
	case HeaderKeyRepoBranchFetched:
		return s.handleRepoBranchFetched(data, header)
	case HeaderKeyRepoFileFetched:
		return s.handleRepoFileFetched(data, header)
	default:
		return errors.New("unknown header value")
	}
}

func (s *server) platform(p string) (codeplatform.CodePlatform, error) {
	v, ok := s.platforms[p]
	if !ok {
		return nil, errors.New("unknown platform: " + p)
	}

	return v, nil
}

func (s *server) handleRepoFetched(data []byte, header map[string]string) error {
	msg := new(msgOfRepoFetched)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	p, err := s.platform(msg.Platform)
	if err != nil {
		return err
	}

	cmd := msg.toCmd()

	return s.service.FetchRepoBranch(p, &cmd)
}

func (s *server) handleRepoBranchFetched(data []byte, header map[string]string) error {
	cmd, platform, err := cmdToFetchRepoFile(data)
	if err != nil {
		return err
	}

	p, err := s.platform(platform)
	if err != nil {
		return err
	}

	return s.service.FetchRepoFile(p, &cmd)
}

func (s *server) handleRepoFileFetched(data []byte, header map[string]string) error {
	cmd, platform, err := cmdToFetchFileContent(data)
	if err != nil {
		return err
	}

	p, err := s.platform(platform)
	if err != nil {
		return err
	}

	return s.service.FetchFileContent(p, &cmd)
}
