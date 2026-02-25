package protocol

// LSP method constants.
const (
	// Lifecycle methods
	MethodInitialize  = "initialize"
	MethodInitialized = "initialized"
	MethodShutdown    = "shutdown"
	MethodExit        = "exit"

	// Document synchronization methods
	MethodDidOpen   = "textDocument/didOpen"
	MethodDidChange = "textDocument/didChange"
	MethodDidClose  = "textDocument/didClose"
	MethodDidSave   = "textDocument/didSave"

	// Notifications from server
	MethodPublishDiagnostics = "textDocument/publishDiagnostics"
)
