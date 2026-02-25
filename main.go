// Package main provides the entry point for the LSP server.
package main

import (
	"os"

	"lsp-server/server"
)

func main() {
	srv := server.New(os.Stdin, os.Stdout)
	if err := srv.Run(); err != nil {
		os.Exit(1)
	}
}
