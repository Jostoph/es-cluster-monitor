package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/Jostoph/es-cluster-monitor/pkg/api"
	"github.com/Jostoph/es-cluster-monitor/pkg/rest/mocks"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestStatusToEnum(t *testing.T) {
	tests := []struct {
		name         string
		stringStatus string
		apiStatus    api.Status
	}{
		{"GREEN", "green", api.Status_GREEN},
		{"YELLOW", "yellow", api.Status_YELLOW},
		{"RED", "red", api.Status_RED},
		{"UNKNOWN", "unknown_status", api.Status_UNKNOWN},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := statusToEnum(test.stringStatus); got != test.apiStatus {
				t.Errorf("expected: %v, got: %v", test.apiStatus, got)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name      string
		stringInt string
		intValue  int64
	}{
		{
			"SUCCESS",
			"42",
			42,
		},
		{
			"FAIL",
			"not_an_int",
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := stringToInt(test.stringInt); got != test.intValue {
				t.Errorf("expected: %d, got: %d", test.intValue, got)
			}
		})
	}
}

func TestESMonitorServer_ReadClusterHealth(t *testing.T) {
	ctx := context.Background()
	httpMockClient := mocks.HTTPClient{}
	server := NewESMonitorServer("", &httpMockClient)

	tests := []struct {
		name      string
		json      string
		wantError bool
		want      *api.ClusterHealthResponse
	}{
		{
			"SUCCESS",
			`{"cluster_name":"docker-cluster","status":"yellow","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":1,"active_shards":1,"relocating_shards":0,"initializing_shards":0,"unassigned_shards":1,"delayed_unassigned_shards":0,"number_of_pending_tasks":0,"number_of_in_flight_fetch":0,"task_max_waiting_in_queue_millis":0,"active_shards_percent_as_number":50.0}`,
			false,
			&api.ClusterHealthResponse{
				ClusterName:                 "docker-cluster",
				Status:                      api.Status_YELLOW,
				NumberOfNodes:               1,
				NumberOfDataNodes:           1,
				ActivePrimaryShards:         1,
				ActiveShards:                1,
				UnassignedShards:            1,
				ActiveShardsPercentAsNumber: 50,
				Timestamp:                   0,
			},
		}, {
			"FAIL",
			"",
			true,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ioutil.NopCloser(bytes.NewReader([]byte(test.json)))

			var fetchError error = nil
			if test.wantError {
				fetchError = errors.New("fetch error")
			}

			httpMockClient.Mock.On("Do", mock.AnythingOfType("*http.Request")).Return(
				func(req *http.Request) *http.Response {
					return &http.Response{
						Body: r,
					}
				},
				func(req *http.Request) error {
					return fetchError
				},
			)

			tm := time.Now().Unix()
			// call service with mocked http client
			got, err := server.ReadClusterHealth(ctx, &api.ClusterHealthRequest{})

			if (err != nil) && !test.wantError {
				t.Errorf("no error expected, got:%v\n", err)
				return
			}

			if err == nil {
				want := test.want
				want.Timestamp = tm
				if !proto.Equal(got, want) {
					t.Errorf("got:\n%+v\n!=\nwant:\n%+v", got, want)
				}
			}
		})

	}
}

func TestESMonitorServer_ReadIndicesInfo(t *testing.T) {
	ctx := context.Background()
	httpMockClient := mocks.HTTPClient{}
	server := NewESMonitorServer("", &httpMockClient)

	tests := []struct {
		name      string
		json      string
		wantError bool
		want      *api.IndicesInfoResponse
	}{
		{
			"SUCCESS",
			`[{"health":"yellow","status":"open","index":"bank","uuid":"IulyimjBSlmck5avEQnHgA","pri":"1","rep":"1","docs.count":"1000","docs.deleted":"0","store.size":"388312","pri.store.size":"388312"}]`,
			false,
			&api.IndicesInfoResponse{
				Indices: []*api.IndexInfo{{
					Health:       api.Status_YELLOW,
					Status:       "open",
					Index:        "bank",
					Uuid:         "IulyimjBSlmck5avEQnHgA",
					Pri:          1,
					Rep:          1,
					DocsCount:    1000,
					StoreSize:    388312,
					PriStoreSize: 388312,
				}},
				Timestamp: 0,
			},
		}, {
			"FAIL",
			"",
			true,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ioutil.NopCloser(bytes.NewReader([]byte(test.json)))

			var fetchError error = nil
			if test.wantError {
				fetchError = errors.New("fetch error")
			}

			httpMockClient.Mock.On("Do", mock.AnythingOfType("*http.Request")).Return(
				func(req *http.Request) *http.Response {
					return &http.Response{
						Body: r,
					}
				},
				func(req *http.Request) error {
					return fetchError
				},
			)

			tm := time.Now().Unix()
			// call service with mocked http client
			got, err := server.ReadIndicesInfo(ctx, &api.IndicesInfoRequest{})

			if (err != nil) && !test.wantError {
				t.Errorf("no error expected, got:%v\n", err)
				return
			}

			if err == nil {
				want := test.want
				want.Timestamp = tm
				if !proto.Equal(got, want) {
					t.Errorf("got:\n%+v\n!=\nwant:\n%+v", got, want)
				}
			}
		})

	}
}
