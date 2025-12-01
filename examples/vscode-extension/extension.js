const { LanguageClient, TransportKind } = require('vscode-languageclient/node');
const vscode = require('vscode');
const path = require('path');

let client;

function activate(context) {
    console.log('Basic LSP Server extension is now active');

    // Get the server path from configuration
    const config = vscode.workspace.getConfiguration('basicLspServer');
    let serverPath = config.get('serverPath');

    // If no path configured, try to find it relative to the extension
    if (!serverPath) {
        serverPath = path.join(__dirname, '..', '..', 'lsp-server');
    }

    // Server options
    const serverOptions = {
        command: serverPath,
        args: [],
        transport: TransportKind.stdio
    };

    // Client options
    const clientOptions = {
        // Register the server for plain text and text documents
        documentSelector: [
            { scheme: 'file', language: 'plaintext' },
            { scheme: 'file', language: 'text' }
        ],
        synchronize: {
            // Notify the server about file changes to files watched by the extension
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*')
        }
    };

    // Create the language client
    client = new LanguageClient(
        'basicLspServer',
        'Basic LSP Server',
        serverOptions,
        clientOptions
    );

    // Start the client (this will also launch the server)
    client.start();

    console.log('Basic LSP Server client started');
}

function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}

module.exports = {
    activate,
    deactivate
};
