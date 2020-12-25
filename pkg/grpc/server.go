package grpc

import (
	"context"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/api"
	"github/Jostoph/es-cluster-monitor/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func NewServer(ctx context.Context, port int, ESAddr string) error {

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// create new grpc server
	server := grpc.NewServer()

	// create new ES Monitor Service Server
	serviceServer := service.NewESMonitorServer(ESAddr)

	// register Service to grpc server
	api.RegisterMonitorServiceServer(server, serviceServer)

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
	return server.Serve(listen)
}
