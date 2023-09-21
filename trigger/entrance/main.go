package main

import (
	"flag"
	"os"

	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-repository-file/trigger/app"
	"github.com/opensourceways/sync-repository-file/trigger/domain/codeplatform"
	"github.com/opensourceways/sync-repository-file/trigger/infrastructure/gitee"
	"github.com/opensourceways/sync-repository-file/trigger/infrastructure/messageimpl"
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
	logrusutil.ComponentInit("sync-repository-file-trigger")
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

	trigger := triggerImpl{
		service: app.NewRepoService(
			messageimpl.NewRepoMessage(&cfg.Message),
		),
		platforms: map[string]codeplatform.CodePlatform{
			giteeImpl.Platform(): giteeImpl,
		},
	}

	trigger.kickoff(cfg.Repos)
}
