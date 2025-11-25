package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// Simple test client to verify LSP server functionality

type Message struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type TestClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr io.ReadCloser
	msgID  int
}

func NewTestClient(serverPath string) (*TestClient, error) {
	cmd := exec.Command(serverPath)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	
	// Start goroutine to display stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "[SERVER] %s\n", scanner.Text())
		}
	}()
	
	return &TestClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
		stderr: stderr,
		msgID:  1,
	}, nil
}

func (c *TestClient) sendMessage(msg Message) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	message := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
	_, err = c.stdin.Write([]byte(message))
	return err
}

func (c *TestClient) readMessage() (*Message, error) {
	// Read headers
	headers := make(map[string]string)
	for {
		line, err := c.stdout.ReadString('\n')
		if err != nil {
			return nil, err
		}
		
		line = line[:len(line)-1]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		
		if line == "" {
			break
		}
		
		var key, value string
		if _, err := fmt.Sscanf(line, "%s %s", &key, &value); err != nil {
			continue
		}
		headers[key[:len(key)-1]] = value
	}
	
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}
	
	var contentLength int
	if _, err := fmt.Sscanf(contentLengthStr, "%d", &contentLength); err != nil {
		return nil, err
	}
	
	content := make([]byte, contentLength)
	if _, err := io.ReadFull(c.stdout, content); err != nil {
		return nil, err
	}
	
	var msg Message
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, err
	}
	
	return &msg, nil
}

func (c *TestClient) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}

func main() {
	serverPath := "./lsp-server"
	if len(os.Args) > 1 {
		serverPath = os.Args[1]
	}
	
	log.Println("Starting test client...")
	log.Printf("Server path: %s\n", serverPath)
	
	client, err := NewTestClient(serverPath)
	if err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer client.Close()
	
	// Test 1: Initialize
	log.Println("\n=== Test 1: Initialize ===")
	initMsg := Message{
		JSONRPC: "2.0",
		ID:      client.msgID,
		Method:  "initialize",
		Params: map[string]interface{}{
			"processId": os.Getpid(),
			"rootUri":   "file:///tmp/test-workspace",
			"capabilities": map[string]interface{}{
				"textDocument": map[string]interface{}{
					"synchronization": map[string]interface{}{
						"didSave": true,
					},
				},
			},
		},
	}
	client.msgID++
	
	if err := client.sendMessage(initMsg); err != nil {
		log.Fatalf("Failed to send initialize: %v", err)
	}
	
	response, err := client.readMessage()
	if err != nil {
		log.Fatalf("Failed to read initialize response: %v", err)
	}
	
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	log.Printf("Initialize response:\n%s\n", responseJSON)
	
	// Test 2: Initialized notification
	log.Println("\n=== Test 2: Initialized Notification ===")
	initializedMsg := Message{
		JSONRPC: "2.0",
		Method:  "initialized",
		Params:  map[string]interface{}{},
	}
	
	if err := client.sendMessage(initializedMsg); err != nil {
		log.Fatalf("Failed to send initialized: %v", err)
	}
	
	// Test 3: textDocument/didOpen
	log.Println("\n=== Test 3: textDocument/didOpen ===")
	didOpenMsg := Message{
		JSONRPC: "2.0",
		Method:  "textDocument/didOpen",
		Params: map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri":        "file:///tmp/test.txt",
				"languageId": "plaintext",
				"version":    1,
				"text":       "Hello, LSP World!",
			},
		},
	}
	
	if err := client.sendMessage(didOpenMsg); err != nil {
		log.Fatalf("Failed to send didOpen: %v", err)
	}
	log.Println("Sent didOpen notification")
	
	// Test 4: textDocument/didChange
	log.Println("\n=== Test 4: textDocument/didChange ===")
	didChangeMsg := Message{
		JSONRPC: "2.0",
		Method:  "textDocument/didChange",
		Params: map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri":     "file:///tmp/test.txt",
				"version": 2,
			},
			"contentChanges": []interface{}{
				map[string]interface{}{
					"text": "Hello, LSP World!\nThis is line 2.",
				},
			},
		},
	}
	
	if err := client.sendMessage(didChangeMsg); err != nil {
		log.Fatalf("Failed to send didChange: %v", err)
	}
	log.Println("Sent didChange notification")
	
	// Test 5: textDocument/didClose
	log.Println("\n=== Test 5: textDocument/didClose ===")
	didCloseMsg := Message{
		JSONRPC: "2.0",
		Method:  "textDocument/didClose",
		Params: map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri": "file:///tmp/test.txt",
			},
		},
	}
	
	if err := client.sendMessage(didCloseMsg); err != nil {
		log.Fatalf("Failed to send didClose: %v", err)
	}
	log.Println("Sent didClose notification")
	
	// Test 6: Shutdown
	log.Println("\n=== Test 6: Shutdown ===")
	shutdownMsg := Message{
		JSONRPC: "2.0",
		ID:      client.msgID,
		Method:  "shutdown",
	}
	client.msgID++
	
	if err := client.sendMessage(shutdownMsg); err != nil {
		log.Fatalf("Failed to send shutdown: %v", err)
	}
	
	response, err = client.readMessage()
	if err != nil {
		log.Fatalf("Failed to read shutdown response: %v", err)
	}
	
	responseJSON, _ = json.MarshalIndent(response, "", "  ")
	log.Printf("Shutdown response:\n%s\n", responseJSON)
	
	// Test 7: Exit
	log.Println("\n=== Test 7: Exit ===")
	exitMsg := Message{
		JSONRPC: "2.0",
		Method:  "exit",
	}
	
	if err := client.sendMessage(exitMsg); err != nil {
		log.Fatalf("Failed to send exit: %v", err)
	}
	log.Println("Sent exit notification")
	
	log.Println("\n=== All tests completed successfully! ===")
}
