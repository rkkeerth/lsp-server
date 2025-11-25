# Basic LSP Client for VS Code

This is an example VS Code extension that connects to the basic LSP server.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Compile the extension:
```bash
npm run compile
```

3. Build the LSP server:
```bash
cd ../..
make build
```

4. Open this directory in VS Code and press F5 to launch the Extension Development Host.

5. In the Extension Development Host, open a `.txt` file to activate the extension.

## Configuration

You can configure the extension in VS Code settings:

- `basicLspServer.trace.server`: Set to "verbose" to see detailed communication logs
- `basicLspServer.serverPath`: Path to the LSP server executable (default: "lsp-server")

## Testing

1. Open a text file in the Extension Development Host
2. Make changes to the file
3. Check the "Output" panel and select "Basic LSP Server" to see server logs
