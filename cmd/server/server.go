package main

import (
	"context"
	"flag"
	"fmt"
	"github/Jostoph/es-cluster-monitor/pkg/grpc"
	"os"
)

func main() {

	// context
	ctx := context.Background()

	// grpc server port
	port := flag.Int("server-port", 9000, "GRPC server port.")

	// Elastic Search Clusters Address
	esAddr := flag.String("es-addr", "http://localhost:9200", "Elastic Search Clusters Address.")
	flag.Parse()

	if err := grpc.NewServer(ctx, *port, *esAddr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
