package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
	"io/ioutil"
	"net/http"
	"time"
)

type ESMonitorServer struct {
	ESAddr string
}

// Return a new Monitor Service Server with ES server address.
func NewESMonitorServer(esAddr string) api.MonitorServiceServer {
	return &ESMonitorServer{ESAddr: esAddr}
}

// ClusterHealth structure to hold an Elastic Search API response for cluster health.
type clusterHealth struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int32   `json:"number_of_nodes"`
	NumberOfDataNodes           int32   `json:"number_of_data_nodes"`
	ActivePrimaryShards         int32   `json:"active_primary_shards"`
	ActiveShards                int32   `json:"active_shards"`
	RelocatingShards            int32   `json:"relocating_shards"`
	InitializingShards          int32   `json:"initializing_shards"`
	UnassignedShards            int32   `json:"unassigned_shards"`
	DelayedUnassignedShards     int32   `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int32   `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int32   `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int32   `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float32 `json:"active_shards_percent_as_number"`
}

// Convert string status to enum.
func statusToEnum(status string) api.ClusterHealthResponse_Status {
	switch status {
	case "green":
		return api.ClusterHealthResponse_GREEN
	case "yellow":
		return api.ClusterHealthResponse_YELLOW
	case "red":
		return api.ClusterHealthResponse_RED
	default:
		return api.ClusterHealthResponse_UNKNOWN
	}
}

// Convert json Cluster ES API health response to proto message.
func jsonClusterToProto(clusterJSON []byte) (*api.ClusterHealthResponse, error) {
	var cluster clusterHealth
	err := json.Unmarshal(clusterJSON, &cluster)
	if err != nil {
		return nil, err
	}

	return &api.ClusterHealthResponse{
		ClusterName:                 cluster.ClusterName,
		Status:                      statusToEnum(cluster.Status),
		TimedOut:                    cluster.TimedOut,
		NumberOfNodes:               cluster.NumberOfNodes,
		NumberOfDataNodes:           cluster.NumberOfDataNodes,
		ActivePrimaryShards:         cluster.ActivePrimaryShards,
		ActiveShards:                cluster.ActiveShards,
		RelocatingShards:            cluster.RelocatingShards,
		InitializingShards:          cluster.InitializingShards,
		UnassignedShards:            cluster.UnassignedShards,
		DelayedUnassignedShards:     cluster.DelayedUnassignedShards,
		NumberOfPendingTasks:        cluster.NumberOfPendingTasks,
		NumberOfInFlightFetch:       cluster.NumberOfInFlightFetch,
		TaskMaxWaitingInQueueMillis: cluster.TaskMaxWaitingInQueueMillis,
		ActiveShardsPercentAsNumber: cluster.ActiveShardsPercentAsNumber,
		Timestamp:                   time.Now().Unix(),
	}, nil
}

func (server *ESMonitorServer) ReadClusterHealth(ctx context.Context, req *api.ClusterHealthRequest) (*api.ClusterHealthResponse, error) {

	// Retrieve Cluster Health as JSON with ES API.
	res, err := http.Get(fmt.Sprintf("%s/_cluster/health?format=JSON", server.ESAddr))
	if err != nil {
		return nil, err
	}

	// Convert JSON response to proto response.
	clusterJSON, _ := ioutil.ReadAll(res.Body)
	proto, err := jsonClusterToProto(clusterJSON)
	if err != nil {
		return nil, err
	}

	return proto, nil
}
