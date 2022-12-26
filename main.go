package main

import (
	"flag"
	"time"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/sirupsen/logrus"

	"project/xihe-statistics/config"
	"project/xihe-statistics/controller"
	"project/xihe-statistics/infrastructure/messages"
	"project/xihe-statistics/infrastructure/pgsql"
	"project/xihe-statistics/server"
)

func main() {
	logrusutil.ComponentInit("xihe-statistics")
	log := logrus.NewEntry(logrus.StandardLogger())

	// cfg
	var cfg string
	flag.StringVar(&cfg, "conf", "./conf/config.yaml", "指定配置文件路径")
	flag.Parse()
	// loading config file
	err := config.Init(cfg)
	if err != nil {
		panic(err)
	}

	// controller
	controller.Init(log)

	// pgsql
	if err := pgsql.Initialize(config.Conf.PGSQL); err != nil {
		logrus.Fatalf("initialize pgsql failed, err:%s", err.Error())
	}

	// init kafka
	if err := messages.Init(config.Conf.GetMQConfig(), log, config.Conf.Topics); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer messages.Exit(log)

	// mq
	go messages.Run(messages.NewHandler(config.Conf, log), log)

	// gin
	server.StartWebServer(config.Conf.HttpPort, time.Duration(config.Conf.Duration))
}
