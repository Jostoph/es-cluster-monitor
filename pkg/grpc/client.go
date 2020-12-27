package grpc

import (
	"context"
	"fmt"
	"github.com/Jostoph/es-cluster-monitor/pkg/api"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"log"
)

func NewClient(ctx context.Context, serverPort int) error {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", serverPort), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// crate grpc client for ES Monitor Service
	monitorService := api.NewMonitorServiceClient(conn)

	resClusterHealth, err := monitorService.ReadClusterHealth(ctx, &api.ClusterHealthRequest{})
	if err != nil {
		return err
	}
	log.Printf("Cluster Health:\n\n%+v\n\n", proto.MarshalTextString(resClusterHealth))

	resIndicesInfo, err := monitorService.ReadIndicesInfo(ctx, &api.IndicesInfoRequest{})
	if err != nil {
		return err
	}
	log.Printf("Indices Info:\n\n%s\n", proto.MarshalTextString(resIndicesInfo))

	return nil
}
