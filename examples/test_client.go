package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Simple LSP client for testing

func main() {
	// Start the LSP server
	cmd := exec.Command("../lsp-server")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdin pipe: %v\n", err)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}

	// Create reader for responses
	reader := bufio.NewReader(stdout)

	// Send initialize request
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"processId": os.Getpid(),
			"rootUri":   "file:///test",
			"capabilities": map[string]interface{}{
				"textDocument": map[string]interface{}{
					"synchronization": map[string]interface{}{
						"dynamicRegistration": true,
					},
				},
			},
		},
	}

	if err := sendMessage(stdin, initRequest); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending initialize: %v\n", err)
		os.Exit(1)
	}

	// Read initialize response
	response, err := readMessage(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading initialize response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Initialize response:")
	printJSON(response)

	// Send initialized notification
	initializedNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
		"params":  map[string]interface{}{},
	}

	if err := sendMessage(stdin, initializedNotif); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending initialized: %v\n", err)
		os.Exit(1)
	}

	// Send textDocument/didOpen notification
	didOpenNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/didOpen",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri":        "file:///test.txt",
				"languageId": "plaintext",
				"version":    1,
				"text":       "Hello, LSP!",
			},
		},
	}

	if err := sendMessage(stdin, didOpenNotif); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending didOpen: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nSent didOpen notification")

	// Send textDocument/didChange notification
	didChangeNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/didChange",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri":     "file:///test.txt",
				"version": 2,
			},
			"contentChanges": []map[string]interface{}{
				{
					"text": "Hello, LSP World!",
				},
			},
		},
	}

	if err := sendMessage(stdin, didChangeNotif); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending didChange: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Sent didChange notification")

	// Send shutdown request
	shutdownRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "shutdown",
	}

	if err := sendMessage(stdin, shutdownRequest); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending shutdown: %v\n", err)
		os.Exit(1)
	}

	// Read shutdown response
	response, err = readMessage(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading shutdown response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nShutdown response:")
	printJSON(response)

	// Send exit notification
	exitNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "exit",
	}

	if err := sendMessage(stdin, exitNotif); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending exit: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nSent exit notification")

	// Wait for server to exit
	if err := cmd.Wait(); err != nil {
		// Exit code 1 is expected after exit notification
		fmt.Fprintf(os.Stderr, "Server exited: %v\n", err)
	}

	fmt.Println("\nTest completed successfully!")
}

func sendMessage(w io.Writer, message interface{}) error {
	content, err := json.Marshal(message)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	_, err = w.Write([]byte(header + string(content)))
	return err
}

func readMessage(r *bufio.Reader) ([]byte, error) {
	// Read headers
	headers := make(map[string]string)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	// Get content length
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length: %s", contentLengthStr)
	}

	// Read content
	content := make([]byte, contentLength)
	_, err = io.ReadFull(r, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func printJSON(data []byte) {
	var pretty interface{}
	if err := json.Unmarshal(data, &pretty); err != nil {
		fmt.Println(string(data))
		return
	}

	formatted, err := json.MarshalIndent(pretty, "", "  ")
	if err != nil {
		fmt.Println(string(data))
		return
	}

	fmt.Println(string(formatted))
}
