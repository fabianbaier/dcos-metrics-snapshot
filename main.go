package main

import (
	"os"

	"github.com/dcos/dcos-go/store"
	"github.com/sirupsen/logrus"

	"github.com/fabianbaier/dcos-metrics-snapshot/config"
	"github.com/fabianbaier/dcos-metrics-snapshot/client"
	"github.com/fabianbaier/dcos-metrics-snapshot/discovery"

	"github.com/fabianbaier/dcos-metrics-snapshot/crawler"
	"time"
	"math"
	"fmt"
)

func main() {
	logrus.Info("Tool to create a snapshot of your cluster metrics environment")

	config, err := config.Parse(os.Args[1:])
	if err != nil {
		logrus.Errorf("Could not parse config: %s", err)
	}

	if config.EnvVerbose() {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("Log Level was set to debug.")
	}

	c, err := client.New(config)
	if err != nil {
		logrus.Errorf("Could not build http client: %s", err)
	}

	logrus.Infof("HTTP Client Timeout set to: %v", config.HTTPClientTimeout())
	logrus.Infof("Secret set to: %v", config.Secret())
	logrus.Infof("Verbose Mode set to: %v", config.EnvVerbose())

	d := discovery.New(c)

	// getting mesos-id and hostname
	agentList, err := d.AgentList()
	if err != nil {
		logrus.Debugf("No AgentList: %s", err)
	}

	logrus.Infof("Received Agentlist: %v", agentList)

	crawler := crawler.New(c, agentList, config)

	fmt.Println("[")
	crawler.Crawl()


	// Basic usage
	s := store.New()
	s.Set("foo", "fooval")
	s.Set("bar", "barval")

	//blocking for a long time, roughly 292.4 years
	<-time.After(time.Duration(math.MaxInt64))
}