package discovery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"github.com/fabianbaier/dcos-metrics-snapshot/client"
)

type Discovery struct {
	client *client.HTTP
}

type MetadataResponse struct {
	ClusterId string `json:"CLUSTER_ID"`
}

type SlavesResponse struct {
	Slaves      []Slaves `json:"slaves"`
}

type Slaves struct {
	Hostname      string `json:"hostname"`
}

type AgentList struct {
	AgentList []Slaves `json:"agent_list"`
}

// New starts a discovery routine to receive metadata for the cluster it is running in
func New(client *client.HTTP) *Discovery {
	d := &Discovery{
		client: client,
	}
	return d
}

func (d *Discovery) AgentList() (AgentList, error) {

	a := AgentList{}
	resp, err := d.client.Read(d.client.Get("http://leader.mesos:5050/slaves"))
	if err != nil {
		logrus.Debugf("[Discovery] Error when reading agent state data: %s", err)
		return AgentList{}, err
	}

	var result SlavesResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		logrus.Errorf("[Discovery] Could not unmarshal agent state response: %s", err)
		return AgentList{}, err
	}
	logrus.Debugf("[Discovery] Received Slaves response: %v", result)

	for _, s := range result.Slaves {
		a.AgentList = append(a.AgentList, s )
	}

	return a, nil
}

