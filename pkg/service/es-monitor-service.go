package service

import (
	"context"
	"github/Jostoph/es-cluster-monitor/pkg/api"
)

type ESMonitorServer struct {}

func (server *ESMonitorServer) ReadHealth(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {

	cluster := &api.GeneralClusterHealthResponse{
		Epoch:               123,
		Timestamp:           "14:14:14",
		Cluster:             "fake-cluster",
		Status:              1,
		NodeTotal:           1,
		NodeData:            1,
		Shards:              1,
		Pri:                 1,
		Relo:                1,
		Init:                1,
		Unassign:            1,
		PendingTasks:        1,
		MaxTaskWaitTime:     123,
		ActiveShardsPercent: 70.0,
	}
	return &api.HealthResponse {
		Clusters: []*api.GeneralClusterHealthResponse{cluster},
	}, nil
}

