package lsp

// LSP method name constants
const (
	MethodInitialize  = "initialize"
	MethodInitialized = "initialized"
	MethodShutdown    = "shutdown"
	MethodExit        = "exit"

	MethodTextDocumentDidOpen           = "textDocument/didOpen"
	MethodTextDocumentDidChange         = "textDocument/didChange"
	MethodTextDocumentDidClose          = "textDocument/didClose"
	MethodTextDocumentPublishDiagnostics = "textDocument/publishDiagnostics"
)
