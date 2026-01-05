package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rkkeerth/lsp-server/server"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting LSP server")

	// Create and start the LSP server
	srv := server.NewServer(logger)
	
	ctx := context.Background()
	if err := srv.Run(ctx, os.Stdin, os.Stdout); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
