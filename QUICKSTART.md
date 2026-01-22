# Quick Start Guide

This guide helps you get the LSP server running in under 5 minutes.

## Step 1: Build the Server

```bash
# Make sure you have Go 1.21+ installed
go version

# Build the server
make build

# Or manually
go build -o lsp-server .
```

## Step 2: Test the Server

Run the server manually to verify it starts:

```bash
./lsp-server
```

You should see:
```
Starting LSP server...
```

The server is now waiting for LSP messages on stdin. Press Ctrl+C to exit.

## Step 3: Configure Your Editor

### VS Code (Quickest)

1. Install the "Custom Local LSP" extension or create a simple extension
2. Add to your settings.json:

```json
{
  "lsp-custom.servers": {
    "example": {
      "command": "/absolute/path/to/lsp-server",
      "fileTypes": ["txt", "md"]
    }
  }
}
```

### Neovim (with nvim-lspconfig)

Add to your config:

```lua
vim.api.nvim_create_autocmd("FileType", {
  pattern = "text",
  callback = function()
    vim.lsp.start({
      name = "example-lsp",
      cmd = {"/absolute/path/to/lsp-server"},
    })
  end,
})
```

### Vim (with vim-lsp)

```vim
if executable('/absolute/path/to/lsp-server')
    au User lsp_setup call lsp#register_server({
        \ 'name': 'example-lsp',
        \ 'cmd': {server_info->['/absolute/path/to/lsp-server']},
        \ 'allowlist': ['text'],
        \ })
endif
```

## Step 4: Try It Out

1. Open a `.txt` file in your configured editor
2. Type `fun` and wait - you should see completion suggestions
3. Type `function` and hover over it - you should see documentation
4. Create a function definition and try go-to-definition

## Example File

Create `test.txt`:

```
function myFunction() {
    var x = 10
    const y = 20
    return x + y
}

// Try hovering over 'function', 'var', 'const'
// Try completion by typing 'fu' or 'va'
// Try go-to-definition on function names
```

## Troubleshooting

**Nothing happens when I type:**
- Check your editor's LSP log (usually in output panel)
- Verify the server path is correct (use absolute paths)
- Make sure the file type is configured

**Server crashes:**
- Check the log file: `./lsp-server 2> server.log`
- Verify dependencies are installed: `go mod verify`

**Completions don't appear:**
- Some editors require specific trigger characters (., :, ()
- Check the editor's completion settings
- Try manual completion trigger (Ctrl+Space in most editors)

## Next Steps

- Read the full README.md for detailed configuration
- Extend the server with language-specific features
- Add more LSP methods (formatting, code actions, etc.)
- Implement proper parsing for your target language
