package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ytake/kfchc/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
)

const (
	version            = 1.0
	connectorStatusURI = "/connectors/%s/status"
)

type (
	RESTClient struct {
		url        *url.URL
		HTTPClient *http.Client
	}
	// ConnectorStatus for kafka connect / connector status
	ConnectorStatus struct {
		RESTClient
		BasicUsername string
		BasicPassword string
	}
	// Requester リクエスト仕様インターフェース
	Requester interface {
		newRequest(ctx context.Context, method string, connect CurrentStatusConnector, body io.Reader) (*http.Request, error)
	}
)

// user agent
var ua = fmt.Sprintf("go-kfchc/%.1f (%s)", version, runtime.Version())

// retryClient internal
func retryClient(logger log.Logger) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.Logger = logger
	return retryClient.StandardClient()
}

// decodeBody internal
func decodeBody(res *http.Response, out interface{}) error {
	b, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(b, out)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
