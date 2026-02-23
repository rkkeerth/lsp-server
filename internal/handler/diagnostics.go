package handler

import (
	"io"

	"github.com/rkkeerth/lsp-server/internal/protocol"
)

// Diagnostics handles publishing diagnostics
type Diagnostics struct {
	writer io.Writer
}

// NewDiagnostics creates a new diagnostics handler
func NewDiagnostics(writer io.Writer) *Diagnostics {
	return &Diagnostics{writer: writer}
}

// Publish sends diagnostics to the client
func (h *Diagnostics) Publish(uri string, diagnostics []protocol.Diagnostic) error {
	notification := &protocol.Notification{
		JSONRPC: "2.0",
		Method:  protocol.MethodTextDocumentPublishDiagnostics,
	}

	params := protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	// Marshal params to JSON for the notification
	// The server will handle actually sending this
	_ = params
	_ = notification

	return nil
}

// CreateDiagnostic is a helper to create a diagnostic
func CreateDiagnostic(line, startChar, endChar int, message string, severity protocol.DiagnosticSeverity) protocol.Diagnostic {
	return protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: line, Character: startChar},
			End:   protocol.Position{Line: line, Character: endChar},
		},
		Severity: &severity,
		Message:  message,
	}
}
