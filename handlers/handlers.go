package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rkkeerth/lsp-server/document"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

// Handler contains all LSP feature handlers
type Handler struct {
	logger          *zap.Logger
	documentManager *document.Manager
	symbolIndex     *SymbolIndex
}

// NewHandler creates a new handler instance
func NewHandler(logger *zap.Logger, docManager *document.Manager) *Handler {
	return &Handler{
		logger:          logger,
		documentManager: docManager,
		symbolIndex:     NewSymbolIndex(),
	}
}

// Hover provides hover information for a symbol at a position
func (h *Handler) Hover(uri protocol.DocumentURI, pos protocol.Position) *protocol.Hover {
	doc, exists := h.documentManager.Get(uri)
	if !exists {
		return nil
	}

	word := doc.GetWordAt(pos)
	if word == "" {
		return nil
	}

	line := doc.GetLineAt(pos)

	// Provide context-aware hover information
	var contents string
	
	// Check for function definitions
	if strings.Contains(line, "func") && strings.Contains(line, word) {
		contents = fmt.Sprintf("**Function**: `%s`\n\nDefined in this document.", word)
	} else if strings.Contains(line, "var") || strings.Contains(line, ":=") {
		contents = fmt.Sprintf("**Variable**: `%s`\n\nDefined in this document.", word)
	} else if strings.Contains(line, "type") && strings.Contains(line, word) {
		contents = fmt.Sprintf("**Type**: `%s`\n\nDefined in this document.", word)
	} else if strings.Contains(line, "const") {
		contents = fmt.Sprintf("**Constant**: `%s`\n\nDefined in this document.", word)
	} else {
		contents = fmt.Sprintf("**Symbol**: `%s`", word)
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: contents,
		},
	}
}

// Definition finds the definition of a symbol at a position
func (h *Handler) Definition(uri protocol.DocumentURI, pos protocol.Position) []protocol.Location {
	doc, exists := h.documentManager.Get(uri)
	if !exists {
		return nil
	}

	word := doc.GetWordAt(pos)
	if word == "" {
		return nil
	}

	h.logger.Debug("Finding definition", zap.String("word", word))

	// Search for the definition in the current document
	locations := h.findSymbolDefinitions(doc, word)

	return locations
}

// References finds all references to a symbol at a position
func (h *Handler) References(uri protocol.DocumentURI, pos protocol.Position) []protocol.Location {
	doc, exists := h.documentManager.Get(uri)
	if !exists {
		return nil
	}

	word := doc.GetWordAt(pos)
	if word == "" {
		return nil
	}

	h.logger.Debug("Finding references", zap.String("word", word))

	// Search for all references in all open documents
	var locations []protocol.Location

	for _, d := range h.documentManager.GetAll() {
		refs := h.findSymbolReferences(d, word)
		locations = append(locations, refs...)
	}

	return locations
}

// DocumentSymbols returns all symbols in a document
func (h *Handler) DocumentSymbols(uri protocol.DocumentURI) []protocol.DocumentSymbol {
	doc, exists := h.documentManager.Get(uri)
	if !exists {
		return nil
	}

	h.logger.Debug("Getting document symbols", zap.String("uri", string(uri)))

	return h.extractDocumentSymbols(doc)
}

// WorkspaceSymbols searches for symbols across the workspace
func (h *Handler) WorkspaceSymbols(query string) []protocol.SymbolInformation {
	h.logger.Debug("Searching workspace symbols", zap.String("query", query))

	var symbols []protocol.SymbolInformation

	for _, doc := range h.documentManager.GetAll() {
		docSymbols := h.extractDocumentSymbols(doc)
		
		for _, sym := range docSymbols {
			// Filter by query
			if query == "" || strings.Contains(strings.ToLower(sym.Name), strings.ToLower(query)) {
				symbols = append(symbols, protocol.SymbolInformation{
					Name: sym.Name,
					Kind: sym.Kind,
					Location: protocol.Location{
						URI:   doc.URI,
						Range: sym.Range,
					},
				})
			}
		}
	}

	return symbols
}

// Completion provides code completion suggestions
func (h *Handler) Completion(uri protocol.DocumentURI, pos protocol.Position) *protocol.CompletionList {
	doc, exists := h.documentManager.Get(uri)
	if !exists {
		return &protocol.CompletionList{
			IsIncomplete: false,
			Items:        []protocol.CompletionItem{},
		}
	}

	h.logger.Debug("Generating completions", zap.String("uri", string(uri)))

	// Get the current line and word prefix
	line := doc.GetLineAt(pos)
	prefix := h.getWordPrefix(line, int(pos.Character))

	var items []protocol.CompletionItem

	// Add common Go keywords
	keywords := []string{
		"func", "var", "const", "type", "struct", "interface",
		"if", "else", "for", "range", "switch", "case", "default",
		"return", "break", "continue", "goto", "defer", "go",
		"package", "import", "map", "chan", "select",
	}

	for _, keyword := range keywords {
		if strings.HasPrefix(keyword, prefix) {
			items = append(items, protocol.CompletionItem{
				Label:  keyword,
				Kind:   protocol.CompletionItemKindKeyword,
				Detail: "keyword",
			})
		}
	}

	// Add symbols from the current document
	symbols := h.extractDocumentSymbols(doc)
	for _, sym := range symbols {
		if strings.HasPrefix(sym.Name, prefix) {
			items = append(items, protocol.CompletionItem{
				Label:  sym.Name,
				Kind:   h.symbolKindToCompletionKind(sym.Kind),
				Detail: sym.Detail,
			})
		}
	}

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}
}

// GetDiagnostics analyzes a document and returns diagnostics
func (h *Handler) GetDiagnostics(uri protocol.DocumentURI) []protocol.Diagnostic {
	doc, exists := h.documentManager.Get(uri)
	if !exists {
		return nil
	}

	var diagnostics []protocol.Diagnostic

	// Simple syntax checking for demonstration
	for lineNum, line := range doc.Lines {
		// Check for unclosed braces
		openBraces := strings.Count(line, "{")
		closeBraces := strings.Count(line, "}")
		
		if openBraces != closeBraces {
			// This is a simple check; real implementation would be more sophisticated
			continue
		}

		// Check for TODO comments
		if strings.Contains(line, "TODO") {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
				Severity: protocol.DiagnosticSeverityHint,
				Message:  "TODO comment found",
				Source:   "lsp-server",
			})
		}

		// Check for FIXME comments
		if strings.Contains(line, "FIXME") {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
				Severity: protocol.DiagnosticSeverityWarning,
				Message:  "FIXME comment found",
				Source:   "lsp-server",
			})
		}
	}

	return diagnostics
}

// findSymbolDefinitions searches for symbol definitions in a document
func (h *Handler) findSymbolDefinitions(doc *document.Document, symbol string) []protocol.Location {
	var locations []protocol.Location

	// Regular expressions for finding definitions
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`func\s+` + regexp.QuoteMeta(symbol) + `\s*\(`),
		regexp.MustCompile(`func\s+\(\w+\s+\*?\w+\)\s+` + regexp.QuoteMeta(symbol) + `\s*\(`),
		regexp.MustCompile(`type\s+` + regexp.QuoteMeta(symbol) + `\s+(struct|interface)`),
		regexp.MustCompile(`var\s+` + regexp.QuoteMeta(symbol) + `\s+`),
		regexp.MustCompile(`const\s+` + regexp.QuoteMeta(symbol) + `\s*=`),
		regexp.MustCompile(regexp.QuoteMeta(symbol) + `\s*:=`),
	}

	for lineNum, line := range doc.Lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				locations = append(locations, protocol.Location{
					URI: doc.URI,
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(lineNum), Character: 0},
						End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
					},
				})
				break
			}
		}
	}

	return locations
}

// findSymbolReferences searches for symbol references in a document
func (h *Handler) findSymbolReferences(doc *document.Document, symbol string) []protocol.Location {
	var locations []protocol.Location

	wordRegex := regexp.MustCompile(`\b` + regexp.QuoteMeta(symbol) + `\b`)

	for lineNum, line := range doc.Lines {
		matches := wordRegex.FindAllStringIndex(line, -1)
		for _, match := range matches {
			locations = append(locations, protocol.Location{
				URI: doc.URI,
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: uint32(match[0])},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(match[1])},
				},
			})
		}
	}

	return locations
}

// extractDocumentSymbols extracts all symbols from a document
func (h *Handler) extractDocumentSymbols(doc *document.Document) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol

	// Regular expressions for different symbol types
	funcRegex := regexp.MustCompile(`func\s+(?:\(\w+\s+\*?\w+\)\s+)?(\w+)\s*\(`)
	typeRegex := regexp.MustCompile(`type\s+(\w+)\s+(struct|interface)`)
	varRegex := regexp.MustCompile(`var\s+(\w+)\s+`)
	constRegex := regexp.MustCompile(`const\s+(\w+)\s*=`)

	for lineNum, line := range doc.Lines {
		// Find functions
		if matches := funcRegex.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, protocol.DocumentSymbol{
				Name:   matches[1],
				Kind:   protocol.SymbolKindFunction,
				Detail: "function",
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
				SelectionRange: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
			})
		}

		// Find types
		if matches := typeRegex.FindStringSubmatch(line); matches != nil {
			kind := protocol.SymbolKindStruct
			if matches[2] == "interface" {
				kind = protocol.SymbolKindInterface
			}
			symbols = append(symbols, protocol.DocumentSymbol{
				Name:   matches[1],
				Kind:   kind,
				Detail: matches[2],
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
				SelectionRange: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
			})
		}

		// Find variables
		if matches := varRegex.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, protocol.DocumentSymbol{
				Name:   matches[1],
				Kind:   protocol.SymbolKindVariable,
				Detail: "variable",
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
				SelectionRange: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
			})
		}

		// Find constants
		if matches := constRegex.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, protocol.DocumentSymbol{
				Name:   matches[1],
				Kind:   protocol.SymbolKindConstant,
				Detail: "constant",
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
				SelectionRange: protocol.Range{
					Start: protocol.Position{Line: uint32(lineNum), Character: 0},
					End:   protocol.Position{Line: uint32(lineNum), Character: uint32(len(line))},
				},
			})
		}
	}

	return symbols
}

// getWordPrefix extracts the word prefix before the cursor position
func (h *Handler) getWordPrefix(line string, pos int) string {
	if pos > len(line) {
		pos = len(line)
	}

	start := pos
	for start > 0 && isWordChar(rune(line[start-1])) {
		start--
	}

	return line[start:pos]
}

// symbolKindToCompletionKind converts a symbol kind to a completion kind
func (h *Handler) symbolKindToCompletionKind(kind protocol.SymbolKind) protocol.CompletionItemKind {
	switch kind {
	case protocol.SymbolKindFunction, protocol.SymbolKindMethod:
		return protocol.CompletionItemKindFunction
	case protocol.SymbolKindVariable:
		return protocol.CompletionItemKindVariable
	case protocol.SymbolKindConstant:
		return protocol.CompletionItemKindConstant
	case protocol.SymbolKindStruct:
		return protocol.CompletionItemKindStruct
	case protocol.SymbolKindInterface:
		return protocol.CompletionItemKindInterface
	default:
		return protocol.CompletionItemKindText
	}
}

// isWordChar determines if a character is part of a word
func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}
