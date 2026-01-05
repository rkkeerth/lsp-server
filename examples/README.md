# LSP Server Examples

This directory contains example files to demonstrate the LSP server's capabilities.

## Testing with sample.go

The `sample.go` file demonstrates various language features that the LSP server can handle:

### Features Demonstrated

1. **Hover Support**
   - Hover over `Calculator` to see type information
   - Hover over method names like `Add`, `Subtract` to see function information
   - Hover over variables to see their declarations

2. **Go to Definition**
   - Click on `NewCalculator()` call and use "Go to Definition" (F12 in VS Code)
   - Jump to method definitions from their usage
   - Navigate to constant and variable definitions

3. **Find References**
   - Right-click on `Calculator` type and find all references
   - See all places where methods are called
   - Find usage of constants like `MaxValue`

4. **Document Symbols**
   - Open the document outline (Ctrl+Shift+O in VS Code)
   - See all functions, types, constants, and variables listed
   - Quick navigation within the file

5. **Code Completion**
   - Type `calc.` and see method suggestions
   - Start typing keywords like `fu` for `func`
   - Get suggestions for defined symbols

6. **Diagnostics**
   - Notice the TODO comment on line 58 - should show as a hint
   - Notice the FIXME comment on line 63 - should show as a warning
   - Real-time feedback as you type

## Manual Testing

### Using with VS Code

1. Build the LSP server:
   ```bash
   cd ..
   ./build.sh
   ```

2. Configure VS Code to use the LSP server (create `.vscode/settings.json` in your workspace)

3. Open `sample.go` in VS Code

4. Try the features:
   - Hover over symbols
   - Use F12 to go to definitions
   - Use Shift+F12 to find references
   - Use Ctrl+Space for completions
   - Check the Problems panel for diagnostics

### Using with Neovim

1. Configure nvim-lspconfig to use the server

2. Open the sample file:
   ```bash
   nvim sample.go
   ```

3. Use LSP keybindings:
   - `K` for hover
   - `gd` for go to definition
   - `gr` for find references
   - `<C-Space>` for completion

## Expected Behaviors

### Hover on line 7 (Calculator type)
```
**Type**: `Calculator`

Defined in this document.
```

### Hover on line 12 (NewCalculator function)
```
**Function**: `NewCalculator`

Defined in this document.
```

### Completion after typing "calc."
Should show:
- Add
- Subtract
- Multiply
- Divide
- Reset
- GetResult

### Document Symbols
Should show:
- Calculator (struct)
- NewCalculator (function)
- Add (function)
- Subtract (function)
- Multiply (function)
- Divide (function)
- Reset (function)
- GetResult (function)
- MaxValue (constant)
- MinValue (constant)
- globalCalculator (variable)
- main (function)

### Diagnostics
- Line 58: "TODO comment found" (Hint)
- Line 63: "FIXME comment found" (Warning)
