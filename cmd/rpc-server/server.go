package main

import (
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
	"github/Jostoph/es-cluster-monitor/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Could not start to listen: %v", err)
	}

	// create new grpc server
	server := grpc.NewServer()

	// create new ES Monitor Service Server
	serviceServer := service.ESMonitorServer{}

	// Register Service to grpc server
	api.RegisterMonitorServiceServer(server, &serviceServer)

	// interruption signal to stop server (^c)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Stopping grpc server...")
			server.GracefulStop()
		}
	}()

	log.Println("Starting grpc server...")

	if err := server.Serve(conn); err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}
}
