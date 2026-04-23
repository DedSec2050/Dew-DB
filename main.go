package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DedSec2050/dew-db/internal/network"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("DEWDB_ADDR")
	if addr == "" {
		addr = ":6379"
	}

	fmt.Printf("Dew-DB server running on %s\n", addr)

	server := network.Server{Addr: addr}
	if err := server.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server stopped with error: %v\n", err)
		os.Exit(1)
	}
}