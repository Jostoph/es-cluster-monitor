package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ESMonitorServer struct{
	ESAddr string
}

// Return a new Monitor Service Server with ES server address.
func NewESMonitorServer(esAddr string) api.MonitorServiceServer {
	return &ESMonitorServer{ESAddr: esAddr}
}

// Cluster structure to hold ES API health response.
type cluster struct {
	Epoch               string `json:"epoch"`
	Timestamp           string `json:"timestamp"`
	Cluster             string `json:"cluster"`
	Status              string `json:"status"`
	NodeTotal           string `json:"node.total"`
	NodeData            string `json:"node.data"`
	Shards              string `json:"shards"`
	Pri                 string `json:"pri"`
	Relo                string `json:"relo"`
	Init                string `json:"init"`
	Unassign            string `json:"unassign"`
	PendingTasks        string `json:"pending_tasks"`
	MaxTaskWaitTime     string `json:"max_task_wait_time"`
	ActiveShardsPercent string `json:"active_shards_percent"`
}

// Convert string status to enum.
func statusToEnum(status string) api.GeneralClusterHealthResponse_Status {
	switch status {
	case "green":
		return api.GeneralClusterHealthResponse_GREEN
	case "yellow":
		return api.GeneralClusterHealthResponse_YELLOW
	case "red":
		return api.GeneralClusterHealthResponse_RED
	default:
		return api.GeneralClusterHealthResponse_UNKNOWN
	}
}

// Converts a string to int32, return 0 if the conversion fails or the value is "-".
func stringToInt32(s string) int32 {

	if s == "-" {
		return 0
	}

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

// Converts a string to float32, return 0 if the conversion fails and removes trailing '%'.
func stringToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

// Convert json Cluster ES API health response to proto message.
func jsonClustersToProto(clustersJSON []byte) (*api.HealthResponse, error) {
	var clusters []cluster
	err := json.Unmarshal(clustersJSON, &clusters)
	if err != nil {
		log.Printf("Conversion error: %s", err)
		return nil, err
	}
	clustersProto := make([]*api.GeneralClusterHealthResponse, 0)
	for _, c := range clusters {
		proto := api.GeneralClusterHealthResponse{
			Epoch:               stringToInt32(c.Epoch),
			Timestamp:           c.Timestamp,
			Cluster:             c.Cluster,
			Status:              statusToEnum(c.Status),
			NodeTotal:           stringToInt32(c.NodeTotal),
			NodeData:            stringToInt32(c.NodeData),
			Shards:              stringToInt32(c.Shards),
			Pri:                 stringToInt32(c.Pri),
			Relo:                stringToInt32(c.Relo),
			Init:                stringToInt32(c.Init),
			Unassign:            stringToInt32(c.Unassign),
			PendingTasks:        stringToInt32(c.PendingTasks),
			MaxTaskWaitTime:     stringToInt32(c.MaxTaskWaitTime),
			ActiveShardsPercent: stringToFloat32(c.ActiveShardsPercent),
		}
		clustersProto = append(clustersProto, &proto)
	}

	return &api.HealthResponse{
		Clusters: clustersProto,
	}, nil
}

func (server *ESMonitorServer) ReadHealth(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {

	res, err := http.Get(fmt.Sprintf("%s/_cat/health?format=JSON", server.ESAddr))
	if err != nil {
		return nil, err
	}

	clustersJSON, _ := ioutil.ReadAll(res.Body)

	resProto, err := jsonClustersToProto(clustersJSON)
	if err != nil {
		return nil, err
	}

	return resProto, nil
}
