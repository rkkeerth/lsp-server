package handlers

import (
	"testing"

	"go.lsp.dev/protocol"
)

func TestNewSymbolIndex(t *testing.T) {
	index := NewSymbolIndex()
	if index == nil {
		t.Fatal("NewSymbolIndex returned nil")
	}
	if index.symbols == nil {
		t.Fatal("Symbol index map is nil")
	}
}

func TestSymbolIndexAdd(t *testing.T) {
	index := NewSymbolIndex()
	
	info := SymbolInfo{
		Name: "testFunc",
		Kind: protocol.SymbolKindFunction,
		Location: protocol.Location{
			URI: "file:///test.go",
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 10},
			},
		},
		Detail: "test function",
	}

	index.Add("testFunc", info)

	symbols := index.Get("testFunc")
	if len(symbols) != 1 {
		t.Errorf("Expected 1 symbol, got %d", len(symbols))
	}
	if symbols[0].Name != "testFunc" {
		t.Errorf("Expected name 'testFunc', got '%s'", symbols[0].Name)
	}
}

func TestSymbolIndexGet(t *testing.T) {
	index := NewSymbolIndex()
	
	info1 := SymbolInfo{Name: "func1", Kind: protocol.SymbolKindFunction}
	info2 := SymbolInfo{Name: "func1", Kind: protocol.SymbolKindFunction}
	
	index.Add("func1", info1)
	index.Add("func1", info2)

	symbols := index.Get("func1")
	if len(symbols) != 2 {
		t.Errorf("Expected 2 symbols, got %d", len(symbols))
	}

	// Test non-existent symbol
	symbols = index.Get("nonexistent")
	if len(symbols) != 0 {
		t.Errorf("Expected 0 symbols for non-existent key, got %d", len(symbols))
	}
}

func TestSymbolIndexClear(t *testing.T) {
	index := NewSymbolIndex()
	
	index.Add("func1", SymbolInfo{Name: "func1"})
	index.Add("func2", SymbolInfo{Name: "func2"})
	
	index.Clear()

	all := index.GetAll()
	if len(all) != 0 {
		t.Errorf("Expected empty index after clear, got %d symbols", len(all))
	}
}

func TestSymbolIndexRemove(t *testing.T) {
	index := NewSymbolIndex()
	uri1 := protocol.DocumentURI("file:///test1.go")
	uri2 := protocol.DocumentURI("file:///test2.go")
	
	info1 := SymbolInfo{
		Name: "func1",
		Location: protocol.Location{URI: uri1},
	}
	info2 := SymbolInfo{
		Name: "func1",
		Location: protocol.Location{URI: uri2},
	}
	
	index.Add("func1", info1)
	index.Add("func1", info2)

	// Remove symbols from uri1
	index.Remove(uri1)

	symbols := index.Get("func1")
	if len(symbols) != 1 {
		t.Errorf("Expected 1 symbol after removal, got %d", len(symbols))
	}
	if symbols[0].Location.URI != uri2 {
		t.Errorf("Expected remaining symbol from uri2")
	}
}

func TestSymbolIndexGetAll(t *testing.T) {
	index := NewSymbolIndex()
	
	index.Add("func1", SymbolInfo{Name: "func1"})
	index.Add("func2", SymbolInfo{Name: "func2"})
	index.Add("func1", SymbolInfo{Name: "func1_v2"})

	all := index.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 unique symbol names, got %d", len(all))
	}

	// Verify it's a copy (modifying returned map shouldn't affect index)
	all["func3"] = []SymbolInfo{{Name: "func3"}}
	
	allAgain := index.GetAll()
	if len(allAgain) != 2 {
		t.Error("GetAll should return a copy, not the original map")
	}
}
