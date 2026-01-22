# VS Code Extension for Example LSP

This directory contains a minimal VS Code extension that connects to the LSP server.

## Setup

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Build the LSP server** (if not already built):
   ```bash
   cd ../..
   make build
   ```

3. **Configure the server path:**
   
   Edit `package.json` and set the default `serverPath` to the absolute path of your LSP server binary, or configure it in VS Code settings after installation.

## Installation

### Option 1: Development Mode

1. Copy this directory to your VS Code extensions folder:
   - Linux/Mac: `~/.vscode/extensions/example-lsp-client-0.1.0/`
   - Windows: `%USERPROFILE%\.vscode\extensions\example-lsp-client-0.1.0\`

2. Restart VS Code

3. Open settings (Ctrl+,) and search for "Example LSP"

4. Set the server path to your LSP server binary

### Option 2: Test in Extension Development Host

1. Open this directory in VS Code
2. Press F5 to launch the Extension Development Host
3. In the new window, open a `.txt` file
4. The LSP features should now be active

## Usage

Once installed and configured:

1. Open any `.txt` file
2. Try the LSP features:
   - **Completion**: Type `fun`, `var`, `const`, etc. and wait for suggestions
   - **Hover**: Hover over keywords to see documentation
   - **Go to Definition**: Right-click on a symbol and select "Go to Definition"

## Configuration

In VS Code settings (`.vscode/settings.json`):

```json
{
  "exampleLsp.serverPath": "/absolute/path/to/lsp-server",
  "exampleLsp.trace.server": "verbose"
}
```

## Commands

- **Example LSP: Restart** - Restarts the LSP server

## Troubleshooting

**Extension not activating:**
- Check VS Code's Output panel → select "Example LSP" from the dropdown
- Verify the server path is correct and the binary exists
- Check that the binary has execute permissions

**No LSP features working:**
- Open the Output panel (View → Output)
- Select "Example LSP" from the dropdown
- Look for error messages

**Server crashes on startup:**
- Run the server manually to see errors: `./lsp-server`
- Check the server logs
- Verify dependencies with: `cd ../.. && go mod verify`

## Extending

To add more file types:

1. Edit `package.json`, add to `activationEvents`:
   ```json
   "onLanguage:markdown"
   ```

2. Edit `extension.js`, add to `documentSelector`:
   ```javascript
   { scheme: 'file', language: 'markdown' }
   ```
