package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcos/dcos-go/dcos/http/transport"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/fabianbaier/dcos-metrics-snapshot/config"
	"io/ioutil"
)

// HTTP is a convenience wrapper around an http.Client
type HTTP struct {
	http.Client
}

// HTTPResult composes the result of an HTTP request
type HTTPResult struct {
	Code int
	Body []byte
}

// Forces urls to use https:// if necessary
type httpsRoundTripper struct {
	delegate http.RoundTripper
	useHTTPS bool
}

func (res *HTTPResult) String() string {
	return fmt.Sprintf("[%d]: %s", res.Code, string(res.Body))
}

// New returns a new http.Client that handles setting the authentication
// header appropriately for the dcos_insights_plunger service account. It also sets
// the url scheme to use http vs. https based on whether or not
// config.FlagCACertFile was set.
func New(cfg *config.Config) (*HTTP, error) {
	roundTripper, err := getTransport(cfg)
	if err != nil {
		return nil, fmt.Errorf("Could not get transport: %s", err)
	}
	if cfg.Secret() != "" {
		secret := cfg.Secret()
		logrus.Debugf("Secret received: %s", secret)

		var config = struct {
			UID           string `json:"uid"`
			Secret        string `json:"private_key"`
			LoginEndpoint string `json:"login_endpoint"`
		}{}

		kByte := []byte(secret)
		if err := json.Unmarshal(kByte, &config); err != nil {
			logrus.Warningf("[Client] Error in secret: %s", err)
			return nil, err
		}

		logrus.Debugf("UID (Service Account): %s", config.UID)
		logrus.Debugf("Secret (Private Key): %s", config.Secret)
		logrus.Debugf("Endpoint (Authentification Endpoint): %s", config.LoginEndpoint)

		roundTripper, err = transport.NewRoundTripper(roundTripper, transport.OptionCredentials(config.UID, config.Secret, config.LoginEndpoint))
		if err != nil {
			return nil, fmt.Errorf("Could not get round tripper with credentials: %s", err)
		}
	}
	client := http.Client{
		Transport: roundTripper,
		Timeout:   cfg.HTTPClientTimeout(),
	}
	return &HTTP{client}, nil
}

// Read reads the result of an HTTP call.
func (h *HTTP) Read(resp *http.Response, err error) (*HTTPResult, error) {
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("Nil response received")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read response body: %s", err)
	}
	return &HTTPResult{
		Code: resp.StatusCode,
		Body: b,
	}, nil
}

func getTransport(cfg *config.Config) (http.RoundTripper, error) {
	transportOptions := []transport.OptionTransportFunc{}
	tr, err := transport.NewTransport(transportOptions...)
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize HTTP transport: %s", err)
	}
	return &httpsRoundTripper{tr, false}, nil
}

func (h *httpsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL
	if h.useHTTPS {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}
	return h.delegate.RoundTrip(req)
}