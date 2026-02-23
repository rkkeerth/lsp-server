package main

import (
	"log"
	"os"

	"github.com/rkkeerth/lsp-server/server"
)

func main() {
	logger := log.New(os.Stderr, "[LSP] ", log.Ldate|log.Ltime|log.Lshortfile)

	srv := server.NewServer(os.Stdin, os.Stdout, logger)
	if err := srv.Run(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
