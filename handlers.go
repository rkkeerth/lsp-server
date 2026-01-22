package main

import (
	"strings"

	"go.lsp.dev/protocol"
)

// generateCompletions generates completion items for the given position
// This is a basic implementation that provides some common keywords and patterns
func (s *Server) generateCompletions(doc *Document, pos protocol.Position) []protocol.CompletionItem {
	// Basic completion items (this can be extended with actual language analysis)
	items := []protocol.CompletionItem{
		{
			Label:  "function",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Function keyword",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "Declares a function",
			},
			InsertText: "function ${1:name}(${2:params}) {\n\t${3:// body}\n}",
		},
		{
			Label:  "if",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "If statement",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "Conditional if statement",
			},
			InsertText: "if ${1:condition} {\n\t${2:// body}\n}",
		},
		{
			Label:  "for",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "For loop",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "For loop statement",
			},
			InsertText: "for ${1:i := 0}; ${2:i < n}; ${3:i++} {\n\t${4:// body}\n}",
		},
		{
			Label:  "const",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Constant declaration",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "Declares a constant",
			},
			InsertText: "const ${1:name} = ${2:value}",
		},
		{
			Label:  "var",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Variable declaration",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "Declares a variable",
			},
			InsertText: "var ${1:name} ${2:type}",
		},
	}

	// Get the current line and word being typed
	lines := strings.Split(doc.Content, "\n")
	if int(pos.Line) < len(lines) {
		line := lines[pos.Line]
		if int(pos.Character) <= len(line) {
			// Extract word at cursor position
			wordStart := pos.Character
			for wordStart > 0 && isIdentifierChar(rune(line[wordStart-1])) {
				wordStart--
			}
			
			prefix := line[wordStart:pos.Character]
			
			// Filter completions based on prefix
			if prefix != "" {
				filtered := []protocol.CompletionItem{}
				for _, item := range items {
					if strings.HasPrefix(strings.ToLower(item.Label), strings.ToLower(prefix)) {
						filtered = append(filtered, item)
					}
				}
				return filtered
			}
		}
	}

	return items
}

// generateHover generates hover information for the given position
func (s *Server) generateHover(doc *Document, pos protocol.Position) *protocol.Hover {
	lines := strings.Split(doc.Content, "\n")
	if int(pos.Line) >= len(lines) {
		return nil
	}

	line := lines[pos.Line]
	if int(pos.Character) >= len(line) {
		return nil
	}

	// Extract word at cursor position
	wordStart := pos.Character
	wordEnd := pos.Character

	for wordStart > 0 && isIdentifierChar(rune(line[wordStart-1])) {
		wordStart--
	}
	for int(wordEnd) < len(line) && isIdentifierChar(rune(line[wordEnd])) {
		wordEnd++
	}

	if wordStart >= wordEnd {
		return nil
	}

	word := line[wordStart:wordEnd]

	// Provide hover information for common keywords
	hoverMap := map[string]string{
		"function": "**function**\n\nDeclares a function in the code.",
		"if":       "**if**\n\nConditional statement that executes code if a condition is true.",
		"for":      "**for**\n\nLoop statement that repeats code.",
		"const":    "**const**\n\nDeclares a constant value that cannot be changed.",
		"var":      "**var**\n\nDeclares a variable.",
		"return":   "**return**\n\nReturns a value from a function.",
		"import":   "**import**\n\nImports a package or module.",
		"package":  "**package**\n\nDeclares the package name for a Go file.",
	}

	if hoverText, ok := hoverMap[word]; ok {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: hoverText,
			},
			Range: &protocol.Range{
				Start: protocol.Position{Line: pos.Line, Character: wordStart},
				End:   protocol.Position{Line: pos.Line, Character: wordEnd},
			},
		}
	}

	// Default hover for any identifier
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "**" + word + "**\n\nIdentifier in the document.",
		},
		Range: &protocol.Range{
			Start: protocol.Position{Line: pos.Line, Character: wordStart},
			End:   protocol.Position{Line: pos.Line, Character: wordEnd},
		},
	}
}

// findDefinition finds the definition location for the symbol at the given position
func (s *Server) findDefinition(doc *Document, pos protocol.Position, uri string) []protocol.Location {
	lines := strings.Split(doc.Content, "\n")
	if int(pos.Line) >= len(lines) {
		return nil
	}

	line := lines[pos.Line]
	if int(pos.Character) >= len(line) {
		return nil
	}

	// Extract word at cursor position
	wordStart := pos.Character
	wordEnd := pos.Character

	for wordStart > 0 && isIdentifierChar(rune(line[wordStart-1])) {
		wordStart--
	}
	for int(wordEnd) < len(line) && isIdentifierChar(rune(line[wordEnd])) {
		wordEnd++
	}

	if wordStart >= wordEnd {
		return nil
	}

	word := line[wordStart:wordEnd]

	// Search for the definition in the document
	// This is a simple implementation that looks for common patterns
	// In a real LSP, this would use proper parsing and symbol resolution
	locations := []protocol.Location{}

	patterns := []string{
		"func " + word,
		"const " + word,
		"var " + word,
		"type " + word,
		word + " :=",
		word + " =",
	}

	for i, docLine := range lines {
		for _, pattern := range patterns {
			if strings.Contains(docLine, pattern) {
				// Found a potential definition
				col := uint32(strings.Index(docLine, word))
				locations = append(locations, protocol.Location{
					URI: protocol.DocumentURI(uri),
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(i), Character: col},
						End:   protocol.Position{Line: uint32(i), Character: col + uint32(len(word))},
					},
				})
			}
		}
	}

	if len(locations) > 0 {
		return locations
	}

	// If no definition found, return the current position
	return []protocol.Location{
		{
			URI: protocol.DocumentURI(uri),
			Range: protocol.Range{
				Start: protocol.Position{Line: pos.Line, Character: wordStart},
				End:   protocol.Position{Line: pos.Line, Character: wordEnd},
			},
		},
	}
}

// isIdentifierChar checks if a character is valid in an identifier
func isIdentifierChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_' || ch == '$'
}
