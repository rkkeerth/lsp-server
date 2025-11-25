package lsp

import (
	"encoding/json"
	"log"
	"os"
)

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(request Request) {
	if s.GetState() != StateUninitialized {
		log.Println("Server already initialized")
		s.sendError(request.ID, InvalidRequest, "Server already initialized")
		return
	}

	s.SetState(StateInitializing)

	// Parse initialize params
	var params InitializeParams
	paramsData, err := json.Marshal(request.Params)
	if err != nil {
		log.Printf("Failed to marshal params: %v", err)
		s.sendError(request.ID, InternalError, "Failed to process parameters")
		return
	}

	if err := json.Unmarshal(paramsData, &params); err != nil {
		log.Printf("Failed to parse initialize params: %v", err)
		s.sendError(request.ID, InvalidParams, "Invalid initialize parameters")
		return
	}

	log.Printf("Client: %s %s", getClientName(params.ClientInfo), getClientVersion(params.ClientInfo))
	log.Printf("Root URI: %s", params.RootURI)

	// Build server capabilities
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: TextDocumentSyncOptions{
				OpenClose: true,
				Change:    Full, // We support full document sync
				Save: &SaveOptions{
					IncludeText: false,
				},
			},
			HoverProvider: false,
			CompletionProvider: &CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{".", ":"},
			},
		},
		ServerInfo: &ServerInfo{
			Name:    "basic-lsp-server",
			Version: "0.1.0",
		},
	}

	s.SetState(StateInitialized)

	if err := s.sendResponse(request.ID, result); err != nil {
		log.Printf("Failed to send initialize response: %v", err)
	}

	log.Println("Server initialized successfully")
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized(request Request) {
	if s.GetState() != StateInitialized {
		log.Println("Received initialized notification but server not initialized")
		return
	}

	log.Println("Client confirmed initialization")
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(request Request) {
	if s.GetState() == StateShuttingDown || s.GetState() == StateShutdown {
		log.Println("Server already shutting down")
		s.sendError(request.ID, InvalidRequest, "Server already shutting down")
		return
	}

	s.SetState(StateShuttingDown)
	log.Println("Server shutting down...")

	// Send null result to indicate success
	if err := s.sendResponse(request.ID, nil); err != nil {
		log.Printf("Failed to send shutdown response: %v", err)
	}
}

// handleExit handles the exit notification
func (s *Server) handleExit(request Request) {
	s.SetState(StateShutdown)
	log.Println("Server exiting...")

	// Exit with appropriate code
	if s.GetState() == StateShuttingDown {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Helper functions
func getClientName(info *ClientInfo) string {
	if info != nil && info.Name != "" {
		return info.Name
	}
	return "unknown"
}

func getClientVersion(info *ClientInfo) string {
	if info != nil && info.Version != "" {
		return info.Version
	}
	return "unknown"
}
