package main

import (
	"context"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
	"google.golang.org/grpc"
	"log"
)

func main() {

	conn, err := grpc.Dial(fmt.Sprintf(":%d", 9000), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server: %s", err)
	}
	defer conn.Close()

	// crate grpc client for ES Monitor Service
	monitorService := api.NewMonitorServiceClient(conn)

	res, err := monitorService.ReadHealth(context.Background(), &api.HealthRequest{})
	if err != nil {
		log.Fatalf("Error while fetching Clusters General Health: %s", err)
	}
	log.Printf("Cluster General Health:\n%+v", res)
}
