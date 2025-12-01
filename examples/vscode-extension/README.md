# VS Code Extension for Basic LSP Server

This is an example VS Code extension that connects to the basic LSP server.

## Setup

1. Build the LSP server first:
   ```bash
   cd ../..
   make build
   ```

2. Install the extension dependencies:
   ```bash
   cd examples/vscode-extension
   npm install
   ```

3. Configure the server path in VS Code settings (optional):
   ```json
   {
     "basicLspServer.serverPath": "/absolute/path/to/lsp-server"
   }
   ```

   If not configured, the extension will look for the binary at `../../lsp-server` (relative to the extension directory).

## Running the Extension

### Option 1: Development Mode

1. Open the `lsp-server` project in VS Code
2. Open the `examples/vscode-extension` folder
3. Press F5 to launch the Extension Development Host
4. In the new VS Code window, open any `.txt` or plaintext file
5. Try hovering over words to see the LSP hover information

### Option 2: Install Locally

1. Package the extension:
   ```bash
   npm install -g vsce
   vsce package
   ```

2. Install the generated `.vsix` file:
   ```bash
   code --install-extension basic-lsp-client-0.1.0.vsix
   ```

3. Reload VS Code and open a plaintext file

## Testing

1. Create a test file: `test.txt`
2. Type some text: `Hello world this is a test`
3. Hover over any word to see hover information
4. The server logs will appear in Output > Basic LSP Server

## Troubleshooting

- Check the Output panel (View > Output) and select "Basic LSP Server" from the dropdown
- Verify the server path is correct in settings
- Make sure the LSP server binary has execute permissions: `chmod +x lsp-server`
- Check that the server works independently: `./lsp-server` (should wait for input)
