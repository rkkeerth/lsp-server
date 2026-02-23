package protocol

// LSP method constants.
const (
	MethodInitialize                   = "initialize"
	MethodInitialized                  = "initialized"
	MethodShutdown                     = "shutdown"
	MethodExit                         = "exit"
	MethodTextDocumentDidOpen          = "textDocument/didOpen"
	MethodTextDocumentDidClose         = "textDocument/didClose"
	MethodTextDocumentDidChange        = "textDocument/didChange"
	MethodTextDocumentCompletion       = "textDocument/completion"
	MethodTextDocumentHover            = "textDocument/hover"
	MethodTextDocumentPublishDiagnostics = "textDocument/publishDiagnostics"
)
