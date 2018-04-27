package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcos/dcos-metrics/producers"
	"github.com/sirupsen/logrus"

	"github.com/fabianbaier/dcos-metrics-snapshot/client"
	"github.com/fabianbaier/dcos-metrics-snapshot/config"
	"github.com/fabianbaier/dcos-metrics-snapshot/discovery"
)

type Crawler struct {
	agentList      discovery.AgentList
	cfg        *config.Config
	client     *client.HTTP
}

func New(client *client.HTTP, agentList discovery.AgentList, cfg *config.Config) *Crawler {

	c := &Crawler{
		agentList:     agentList,
		cfg:       cfg,
		client:    client,
	}
	return c
}

// Start initializes scrapers and pushers
func (c *Crawler) Crawl() error {
	newMetricScraper(c)
	return nil
}

type metricCrawler struct {
	c *Crawler
}

// NewMetricScraper initializes a new metric scraper for a specific agent.
// It scrapes node metrics as well as container metrics.
func newMetricScraper(c *Crawler) *metricCrawler {
	mc := &metricCrawler{
		c: c,
	}


	go func() {
			mc.crawl()
	}()
	return mc
}

func (m *metricCrawler) crawl() {
	go m.fetchAgent()
}

func (m *metricCrawler) fetchAgent() error {
	for _, hostname := range m.c.agentList.AgentList {
		go m.fetchContainerMetrics(hostname)
	}
	return nil
}

func (m *metricCrawler) fetchContainerMetrics(hostname discovery.Slaves) error {
	containerListURL := fmt.Sprintf("http://%s:61001/system/v1/metrics/v0/containers", hostname.Hostname)

	logrus.Debugf("Started container list scraper for URL: %s", containerListURL)

	resp, err := m.c.client.Read(m.c.client.Get(containerListURL))
	if err != nil {
		logrus.Debugf("[Crawler] Could not pull container list: %s", err)
		return err
	}
	if resp.Code != http.StatusOK {
		logrus.Debugf("Polling container list returned %s", resp)
		return err
	}

	result := []string{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		logrus.Debugf("Could not unmarshal dcos-metric api response: %s", err)
		return err
	}

	logrus.Debugf("Received list of %d containers @: %s", len(result), hostname.Hostname)
	for _, c := range result {
		c := c
		h := hostname.Hostname
		// getting running container metrics
		go func() {
			err := m.getContainerMetric(h, c)
			if err != nil {
				logrus.Debugf("failed starting fetching container metrics: %s", err)
			}
		}()
	}
	return nil
}

// getContainerMetric is the scraper method that connects to the agents DC/OS metrics api container endpoint
// and receives metrics of a container id.
func (m *metricCrawler) getContainerMetric(h, c string) error {

	containerMetricURL := fmt.Sprintf("http://%s:61001/system/v1/metrics/v0/containers/%s", h, c)

	logrus.Debugf("Started container metric scrapper for Container URL: %s", containerMetricURL)

	resp, err := m.c.client.Read(m.c.client.Get(containerMetricURL))
	if err != nil {
		logrus.Debugf("Could not pull container metrics: %s", err)
		return err
	}

	if resp.Code == http.StatusNoContent {
		logrus.Debugf("Found pod: %s (Code: %s)", c, resp)
		// todo: Pull pod metrics
	}

	if resp.Code != http.StatusOK {
		logrus.Debugf("Polling container metric returned %s", resp)
		return err
	}

	result := &producers.MetricsMessage{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		logrus.Debugf("Could not unmarshal adminrouter dcos-metric api response: %s", err)
		return err
	}

	e, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(e),",")


	return nil
}