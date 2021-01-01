package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Jostoph/es-cluster-monitor/pkg/grpc"
	"os"
)

func main() {

	// context
	ctx := context.Background()

	// grpc server port
	port := flag.Int("server-port", 9000, "GRPC server port.")

	// elastic search clusters address
	esAddr := flag.String("es-addr", "http://localhost:9200", "Elastic Search Clusters Address.")
	flag.Parse()

	server := grpc.NewServer(ctx, *port, *esAddr)

	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
