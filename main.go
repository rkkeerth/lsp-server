package main

import (
	"context"
	"log"
	"os"

	"go.lsp.dev/jsonrpc2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting LSP server...")

	// Create the LSP server
	server := NewServer()

	// Create a JSON-RPC 2.0 connection over stdio
	ctx := context.Background()
	stream := jsonrpc2.NewStream(stdrwc{})

	conn := jsonrpc2.NewConn(stream)
	server.conn = conn

	// Set the handler for incoming requests and notifications
	conn.Go(ctx, jsonrpc2.HandlerFunc(server.Handle))

	// Wait for the connection to close
	<-conn.Done()

	log.Println("LSP server stopped")
}

// stdrwc implements io.ReadWriteCloser for stdin/stdout
type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
