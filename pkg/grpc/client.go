package grpc

import (
	"context"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
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

	res, err := monitorService.ReadHealth(ctx, &api.HealthRequest{})
	if err != nil {
		return err
	}
	log.Printf("Cluster General Health:\n%+v", res)
	return nil
}