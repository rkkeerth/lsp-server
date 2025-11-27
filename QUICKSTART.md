# Quick Start Guide

Get started with the Basic LSP Server in 2 minutes!

## Prerequisites

- Python 3.8 or higher
- Any LSP-compatible editor (VS Code, Neovim, Emacs, etc.)

## Quick Test

Run the automated test:

```bash
./run_test.sh
```

This will:
- Verify Python installation
- Test the server with example files
- Show diagnostic output

## Manual Test

1. **Start testing with a single command:**
   ```bash
   python3 test_client.py examples/test.txt
   ```

2. **Test with your own file:**
   ```bash
   python3 test_client.py /path/to/your/file.txt
   ```

## What Gets Detected?

The server will flag:

- ✓ `TODO` comments → Information
- ⚠ `FIXME` comments → Warning
- ⚠ Lines over 120 characters → Warning
- ✓ Duplicate consecutive lines → Information

## Example Output

```
[INFO] Line 5, Col 1: TODO found: Consider addressing this item
[WARNING] Line 8, Col 1: FIXME found: This requires immediate attention
[WARNING] Line 10, Col 121: Line too long (136 > 120 characters)
```

## Connect to Your Editor

### VS Code
1. Install an LSP client extension
2. Configure server path: `/path/to/server.py`
3. Restart VS Code

### Neovim
Add to your config:
```lua
vim.lsp.start({
  name = 'basic-lsp-server',
  cmd = {'python3', '/path/to/server.py'},
})
```

### Emacs
Add to your config:
```elisp
(lsp-register-client
 (make-lsp-client 
   :new-connection (lsp-stdio-connection '("python3" "/path/to/server.py"))
   :major-modes '(text-mode)
   :server-id 'basic-lsp-server))
```

## Troubleshooting

**Server not responding?**
```bash
tail -f /tmp/lsp-server.log
```

**Python not found?**
```bash
python3 --version
```

**Need help?**
See [README.md](README.md) for detailed documentation.

## Next Steps

1. Try editing `examples/test.txt` and watch diagnostics update
2. Add your own diagnostic rules in `server.py`
3. Connect the server to your favorite editor

Happy coding! 🚀
