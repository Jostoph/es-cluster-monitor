package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
	"io/ioutil"
	"net/http"
	"strconv"
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

// IndexInfo structure to hold an Elastic Search API response for indices info.
type indexInfo struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	Uuid         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

// Convert string status to enum.
func statusToEnum(status string) api.Status {
	switch status {
	case "green":
		return api.Status_GREEN
	case "yellow":
		return api.Status_YELLOW
	case "red":
		return api.Status_RED
	default:
		return api.Status_UNKNOWN
	}
}

// Convert string to int32. Return 0 if the conversion fails.
func stringToInt32(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}

	return int32(i)
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

func jsonIndicesToProto(indicesJSON []byte) (*api.IndicesInfoResponse, error) {
	var indices []indexInfo
	err := json.Unmarshal(indicesJSON, &indices)
	if err != nil {
		return nil, err
	}

	indicesProto := make([]*api.IndexInfo, 0)
	for _, idxInfo := range indices {
		proto := api.IndexInfo{
			Health:       statusToEnum(idxInfo.Status),
			Status:       idxInfo.Status,
			Index:        idxInfo.Index,
			Uuid:         idxInfo.Uuid,
			Pri:          stringToInt32(idxInfo.Pri),
			Rep:          stringToInt32(idxInfo.Rep),
			DocsCount:    stringToInt32(idxInfo.DocsCount),
			DocsDeleted:  stringToInt32(idxInfo.DocsDeleted),
			StoreSize:    stringToInt32(idxInfo.StoreSize),
			PriStoreSize: stringToInt32(idxInfo.PriStoreSize),
		}
		indicesProto = append(indicesProto, &proto)
	}

	return &api.IndicesInfoResponse{
		Indices:   indicesProto,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (server *ESMonitorServer) ReadClusterHealth(ctx context.Context, req *api.ClusterHealthRequest) (*api.ClusterHealthResponse, error) {

	// Retrieve Cluster Health as JSON with ES API.
	res, err := http.Get(fmt.Sprintf("%s/_cluster/health?format=JSON", server.ESAddr))
	if err != nil {
		return nil, err
	}

	// Convert JSON response to proto response.
	clusterJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	proto, err := jsonClusterToProto(clusterJSON)
	if err != nil {
		return nil, err
	}

	return proto, nil
}

func (server *ESMonitorServer) ReadIndicesInfo(ctx context.Context, req *api.IndicesInfoRequest) (*api.IndicesInfoResponse, error) {

	// Retrieve Indices Info as JSON with ES API.
	res, err := http.Get(fmt.Sprintf("%s/_cat/indices?v&bytes=b&format=JSON", server.ESAddr))
	if err != nil {
		return nil, err
	}

	// Convert JSON response to proto response.
	indicesJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	proto, err := jsonIndicesToProto(indicesJSON)
	if err != nil {
		return nil, err
	}

	return proto, nil
}
