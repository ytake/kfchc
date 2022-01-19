package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	responseOK    = `{"name":"file_sink","connector":{"state":"RUNNING","worker_id":"kafka-connect:8083"},"tasks":[{"id":0,"state":"RUNNING","worker_id":"kafka-connect:8083"}],"type":"sink"}`
	responseError = `{"name":"file_json_sink","connector":{"state":"RUNNING","worker_id":"kafka-connect:8083"},"tasks":[{"id":0,"state":"FAILED","worker_id":"kafka-connect:8083","trace":"error."}],"type":"sink"}`
)

type ts struct {
	name                string
	mockResponseBody    string
	expectedMethod      string
	expectedRequestPath string
	expectedErrMessage  string
}

func TestConnectorStatus_Get(t *testing.T) {
	var tt []ts
	tt = append(tt, ts{
		name:                "sample1",
		expectedMethod:      http.MethodGet,
		expectedRequestPath: "/connectors/sample1/status",
		mockResponseBody:    responseOK,
	})
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != tc.expectedMethod {
					t.Fatalf("request method wrong. want=%s, got=%s", tc.expectedMethod, req.Method)
				}
				if req.URL.Path != tc.expectedRequestPath {
					t.Fatalf("request path wrong. want=%s, got=%s", tc.expectedRequestPath, req.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				bodyBytes, _ := ioutil.ReadAll(strings.NewReader(tc.mockResponseBody))
				w.Write(bodyBytes)
			}))
			defer server.Close()
			serverURL, _ := url.Parse(server.URL)
			cs := &ConnectorStatus{
				RESTClient: RESTClient{
					url:        serverURL,
					HTTPClient: server.Client(),
				},
			}
			csc := &CurrentStatus{ConnectorName: "sample1"}
			r := cs.Get(csc, context.TODO())
			rs := <-r
			if rs.ConnectorStatus.Name != "file_sink" {
				t.Error("connector name wrong.")
			}
			if rs.ConnectorStatus.Connector.IsFailed() {
				t.Error("connector status wrong.")
			}
		})
	}
}

func TestConnectorStatusFailed_Get(t *testing.T) {
	var tt []ts
	tt = append(tt, ts{
		name:                "sample2",
		expectedMethod:      http.MethodGet,
		expectedRequestPath: "/connectors/sample2/status",
		mockResponseBody:    responseError,
	})
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != tc.expectedMethod {
					t.Fatalf("request method wrong. want=%s, got=%s", tc.expectedMethod, req.Method)
				}
				if req.URL.Path != tc.expectedRequestPath {
					t.Fatalf("request path wrong. want=%s, got=%s", tc.expectedRequestPath, req.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				bodyBytes, _ := ioutil.ReadAll(strings.NewReader(tc.mockResponseBody))
				w.Write(bodyBytes)
			}))
			defer server.Close()
			serverURL, _ := url.Parse(server.URL)
			cs := &ConnectorStatus{
				RESTClient: RESTClient{
					url:        serverURL,
					HTTPClient: server.Client(),
				},
			}
			csc := &CurrentStatus{ConnectorName: "sample2"}
			r := cs.Get(csc, context.TODO())
			rs := <-r
			if rs.ConnectorStatus.Name != "file_json_sink" {
				t.Error("connector name wrong.")
			}
			if rs.ConnectorStatus.Connector.IsFailed() {
				t.Error("connector status wrong.")
			}
			if len(rs.ConnectorStatus.Tasks) == 0 {
				t.Error("connector tasks wrong.")
			}
			if !rs.ConnectorStatus.Tasks[0].IsFailed() {
				t.Error("task status wrong.")
			}
		})
	}
}
