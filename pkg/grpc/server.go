package grpc

import (
	"context"
	"fmt"
	"github.com/Jostoph/es-cluster-monitor/pkg/api"
	"github.com/Jostoph/es-cluster-monitor/pkg/logger"
	"github.com/Jostoph/es-cluster-monitor/pkg/rest"
	"github.com/Jostoph/es-cluster-monitor/pkg/service"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

type Server struct {
	ctx  context.Context
	port int
	addr string
}

func NewServer(ctx context.Context, port int, ESAddr string) *Server {
	return &Server{
		ctx:  ctx,
		port: port,
		addr: ESAddr,
	}
}

func (server *Server) Run() error {

	// init logger
	logger.Init()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		return err
	}

	// create new grpc server
	serv := grpc.NewServer()

	// create new ES Monitor Service Server
	serviceServer := service.NewESMonitorServer(server.addr, rest.NewClient())

	// register Service to grpc server
	api.RegisterMonitorServiceServer(serv, serviceServer)

	// interruption signal to stop server (^c)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Log.Warn("Stopping grpc server.")
			serv.GracefulStop()
		}
	}()

	logger.Log.Info("Starting grpc server.")
	return serv.Serve(listen)
}
