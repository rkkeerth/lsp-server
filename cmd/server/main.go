// Package main is the entry point for the LSP server.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/example/lsp-server/internal/lsp"
)

const (
	version = "0.1.0"
)

func main() {
	// Parse command-line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help message")
		tcpMode     = flag.Bool("tcp", false, "Run server in TCP mode")
		tcpAddr     = flag.String("addr", "127.0.0.1:7998", "TCP address to listen on (when -tcp is used)")
		logFile     = flag.String("log", "", "Log file path (default: stderr)")
	)
	flag.Parse()

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("lsp-server version %s\n", version)
		os.Exit(0)
	}

	// Set up logging
	logger := setupLogger(*logFile)

	logger.Printf("Starting LSP server v%s", version)

	if *tcpMode {
		runTCPServer(*tcpAddr, logger)
	} else {
		runStdioServer(logger)
	}
}

// printUsage prints the help message.
func printUsage() {
	fmt.Println("LSP Server - A Language Server Protocol implementation in Go")
	fmt.Println()
	fmt.Println("Usage: lsp-server [options]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("By default, the server runs in stdio mode, reading from stdin and writing to stdout.")
	fmt.Println("Use -tcp flag to run in TCP mode for debugging or testing.")
}

// setupLogger configures the logger.
func setupLogger(logFile string) *log.Logger {
	var writer io.Writer = os.Stderr

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		writer = file
	}

	return log.New(writer, "[LSP] ", log.LstdFlags|log.Lshortfile)
}

// runStdioServer runs the server using stdin/stdout.
func runStdioServer(logger *log.Logger) {
	logger.Println("Running in stdio mode")

	server := lsp.NewServer(os.Stdin, os.Stdout, logger)
	if err := server.Run(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

// runTCPServer runs the server in TCP mode.
func runTCPServer(addr string, logger *log.Logger) {
	logger.Printf("Running in TCP mode on %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer listener.Close()

	logger.Printf("Listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("Failed to accept connection: %v", err)
			continue
		}

		logger.Printf("New connection from %s", conn.RemoteAddr())

		// Handle each connection in a separate goroutine
		go func(c net.Conn) {
			defer c.Close()

			server := lsp.NewServer(c, c, logger)
			if err := server.Run(); err != nil {
				logger.Printf("Server error: %v", err)
			}

			logger.Printf("Connection closed: %s", c.RemoteAddr())
		}(conn)
	}
}
