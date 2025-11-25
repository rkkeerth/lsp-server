package main

import (
	"log"
	"os"

	"github.com/rkkeerth/lsp-server/internal/lsp"
)

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stderr) // LSP uses stdout for communication

	log.Println("Starting LSP Server...")

	// Create and start the server
	// LSP communication happens via stdin/stdout
	server := lsp.NewServer(os.Stdin, os.Stdout)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
