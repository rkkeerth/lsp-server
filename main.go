package main

import (
	"context"
	"log"
	"os"

	"go.lsp.dev/jsonrpc2"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Starting LSP server...")

	// Create a new LSP server
	server := NewServer()

	// Create JSON-RPC stream over stdin/stdout
	stream := jsonrpc2.NewStream(&StdioReadWriteCloser{
		in:  os.Stdin,
		out: os.Stdout,
	})

	// Create connection with the stream
	conn := jsonrpc2.NewConn(stream)

	// Set the connection in the server
	server.conn = conn

	// Start handling requests
	ctx := context.Background()
	handler := jsonrpc2.HandlerFunc(server.Handle)

	// Run the connection with our handler
	<-conn.Go(ctx, handler).Done()

	log.Println("LSP server shutting down...")
}

// StdioReadWriteCloser wraps stdin/stdout for JSON-RPC communication
type StdioReadWriteCloser struct {
	in  *os.File
	out *os.File
}

func (s *StdioReadWriteCloser) Read(p []byte) (int, error) {
	return s.in.Read(p)
}

func (s *StdioReadWriteCloser) Write(p []byte) (int, error) {
	return s.out.Write(p)
}

func (s *StdioReadWriteCloser) Close() error {
	if err := s.in.Close(); err != nil {
		return err
	}
	return s.out.Close()
}
