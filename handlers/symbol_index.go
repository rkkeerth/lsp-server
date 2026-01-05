package handlers

import (
	"sync"

	"go.lsp.dev/protocol"
)

// SymbolInfo contains information about a symbol
type SymbolInfo struct {
	Name     string
	Kind     protocol.SymbolKind
	Location protocol.Location
	Detail   string
}

// SymbolIndex maintains an index of symbols for quick lookup
type SymbolIndex struct {
	symbols map[string][]SymbolInfo
	mu      sync.RWMutex
}

// NewSymbolIndex creates a new symbol index
func NewSymbolIndex() *SymbolIndex {
	return &SymbolIndex{
		symbols: make(map[string][]SymbolInfo),
	}
}

// Add adds a symbol to the index
func (si *SymbolIndex) Add(name string, info SymbolInfo) {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.symbols[name] = append(si.symbols[name], info)
}

// Get retrieves all symbols with the given name
func (si *SymbolIndex) Get(name string) []SymbolInfo {
	si.mu.RLock()
	defer si.mu.RUnlock()

	return si.symbols[name]
}

// Clear removes all symbols from the index
func (si *SymbolIndex) Clear() {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.symbols = make(map[string][]SymbolInfo)
}

// Remove removes all symbols associated with a specific URI
func (si *SymbolIndex) Remove(uri protocol.DocumentURI) {
	si.mu.Lock()
	defer si.mu.Unlock()

	for name, infos := range si.symbols {
		filtered := make([]SymbolInfo, 0)
		for _, info := range infos {
			if info.Location.URI != uri {
				filtered = append(filtered, info)
			}
		}
		if len(filtered) > 0 {
			si.symbols[name] = filtered
		} else {
			delete(si.symbols, name)
		}
	}
}

// GetAll returns all symbols in the index
func (si *SymbolIndex) GetAll() map[string][]SymbolInfo {
	si.mu.RLock()
	defer si.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string][]SymbolInfo)
	for k, v := range si.symbols {
		result[k] = append([]SymbolInfo{}, v...)
	}
	return result
}
