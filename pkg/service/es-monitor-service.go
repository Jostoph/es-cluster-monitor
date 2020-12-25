package service

import (
	"context"
	"github/Jostoph/es-cluster-monitor/pkg/api"
)

type ESMonitorServer struct {}

func (server *ESMonitorServer) ReadHealth(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {

	cluster := &api.GeneralClusterHealthResponse{
		Epoch:               0,
		Timestamp:           "14:14:14",
		Cluster:             "fake-cluster",
		Status:              0,
		NodeTotal:           0,
		NodeData:            0,
		Shards:              0,
		Pri:                 0,
		Relo:                0,
		Init:                0,
		Unassign:            0,
		PendingTasks:        0,
		MaxTaskWaitTime:     0,
		ActiveShardsPercent: 50.0,
	}
	return &api.HealthResponse {
		Clusters: []*api.GeneralClusterHealthResponse{cluster},
	}, nil
}

