package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-repository-file/server/app"
	"github.com/opensourceways/sync-repository-file/server/domain/codeplatform"
	"github.com/opensourceways/sync-repository-file/server/infrastructure/gitee"
	"github.com/opensourceways/sync-repository-file/server/infrastructure/messageimpl"
	"github.com/opensourceways/sync-repository-file/server/infrastructure/repositoryimpl"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false, "whether to enable debug model.",
	)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("sync-repository-file-server")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg, err := loadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Errorf("load config, err:%s", err.Error())

		return
	}

	// mq
	if err = kafka.Init(&cfg.Kafka, log, nil, ""); err != nil {
		logrus.Errorf("initialize mq failed, err:%v", err)

		return
	}

	defer kafka.Exit()

	// service
	giteeImpl := gitee.NewGiteePlatform(&cfg.Gitee)

	s := server{
		service: app.NewRepoFileService(
			repositoryimpl.NewRepoFileImpl(&cfg.Repository),
			messageimpl.NewRepoFileMessage(&cfg.Message),
		),
		platforms: map[string]codeplatform.CodePlatform{
			giteeImpl.Platform(): giteeImpl,
		},
	}

	// run
	run(&s, cfg)
}

func run(s *server, cfg *Config) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")

			return

		case <-sig:
			logrus.Info("receive exit signal")
			done()
			called = true

			return
		}
	}(ctx)

	if err := s.run(ctx, cfg); err != nil {
		logrus.Errorf("server exited, err:%s", err.Error())
	}
}
