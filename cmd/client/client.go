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
	serverPort := flag.Int("server-port", 9000, "GRPC server port.")
	flag.Parse()

	if err := grpc.NewClient(ctx, *serverPort); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
