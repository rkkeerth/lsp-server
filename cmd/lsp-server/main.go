package main

import (
	"fmt"
	"os"

	"github.com/rkkeerth/lsp-server/internal/jsonrpc"
	"github.com/rkkeerth/lsp-server/internal/server"
)

func main() {
	// Create LSP server
	lspServer := server.NewServer()

	// Create JSON-RPC handler
	rpc := jsonrpc.NewRPC(os.Stdin, os.Stdout, lspServer)

	// Start the server
	if err := rpc.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
