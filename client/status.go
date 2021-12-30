package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ytake/kfchc/log"
	"github.com/ytake/kfchc/payload"
	"io"
	"net/http"
	"net/url"
	"path"
)

func NewKafkaConnect(connectServer string, logger log.Logger) (*ConnectorStatus, error) {
	u, err := url.Parse(connectServer)
	if err != nil {
		return nil, err
	}
	return &ConnectorStatus{
		RESTClient: RESTClient{
			url:        u,
			HTTPClient: retryClient(logger),
		},
	}, nil
}

func (cs *ConnectorStatus) newRequest(ctx context.Context, method string, connect CurrentStatusConnector, body io.Reader) (*http.Request, error) {
	u := *cs.url
	u.Path = path.Join(connect.DetectPath())
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if cs.BasicPassword != "" && cs.BasicUsername != "" {
		req.SetBasicAuth(cs.BasicUsername, cs.BasicPassword)
	}
	req.Header.Set("User-Agent", ua)
	return req, nil
}

// Get fo kafka connector status
func (cs *ConnectorStatus) Get(connect CurrentStatusConnector, ctx context.Context) <-chan payload.ResultConnectorStatus {
	out := make(chan payload.ResultConnectorStatus)
	go func() {
		defer close(out)
		var result payload.ResultConnectorStatus
		req, err := cs.newRequest(ctx, http.MethodGet, connect, bytes.NewBuffer([]byte{}))
		if err != nil {
			result.Err = err
			out <- result
			return
		}
		res, err := cs.HTTPClient.Do(req)
		defer res.Body.Close()
		if err != nil {
			result.Err = err
			out <- result
			return
		}
		if res.StatusCode == http.StatusNotFound {
			var errorDef payload.ConnectorNotFound
			if err := decodeBody(res, &errorDef); err != nil {
				result.Err = err
				out <- result
				return
			}
			result.ConnectorNotFound = errorDef
			out <- result
			return
		}
		if res.StatusCode == http.StatusOK {
			var definition payload.ConnectorStatus
			if err := decodeBody(res, &definition); err != nil {
				result.Err = err
				out <- result
				return
			}
			result.ConnectorStatus = definition
			out <- result
			return
		}
		result.Err = errors.New("error")
		out <- result
	}()
	return out
}

type CurrentStatusConnector interface {
	DetectPath() string
}

type CurrentStatus struct {
	ConnectorName string
}

func (cs *CurrentStatus) DetectPath() string {
	return fmt.Sprintf(connectorStatusURI, cs.ConnectorName)
}
