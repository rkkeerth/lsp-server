package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rkkeerth/lsp-server/internal/server"
)

var version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("lsp-server version %s\n", version)
		os.Exit(0)
	}

	// Create logger that writes to stderr (stdout is for LSP messages)
	logger := log.New(os.Stderr, "[lsp-server] ", log.LstdFlags)
	logger.Println("Starting LSP server...")

	// Create and run the server
	srv := server.New(os.Stdin, os.Stdout, logger)
	if err := srv.Run(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
