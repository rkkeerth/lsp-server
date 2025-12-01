package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/rkkeerth/lsp-server/server"
	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	// Create a logger that writes to stderr (stdout is used for LSP communication)
	logger := log.New(os.Stderr, "[LSP] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting LSP server...")

	// Create the LSP server
	srv := server.NewServer(logger)

	// Create a JSON-RPC 2.0 connection over stdin/stdout
	conn := jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(srv.Handle),
	)

	logger.Println("LSP server is ready and listening on stdin/stdout")

	// Wait for the connection to close
	<-conn.DisconnectNotify()
	logger.Println("LSP server connection closed")
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

var _ io.ReadWriteCloser = stdrwc{}
