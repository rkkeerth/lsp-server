const vscode = require('vscode');
const { LanguageClient, TransportKind } = require('vscode-languageclient/node');

let client;

/**
 * Activates the extension
 */
function activate(context) {
    console.log('Example LSP client is now active');

    // Get configuration
    const config = vscode.workspace.getConfiguration('exampleLsp');
    const serverPath = config.get('serverPath') || 'lsp-server';
    const trace = config.get('trace.server') || 'off';

    // Server options
    const serverOptions = {
        command: serverPath,
        args: [],
        transport: TransportKind.stdio
    };

    // Client options
    const clientOptions = {
        // Register the server for plain text documents
        documentSelector: [
            { scheme: 'file', language: 'plaintext' },
            { scheme: 'file', language: 'text' }
        ],
        synchronize: {
            // Synchronize the file configuration section to the server
            configurationSection: 'exampleLsp',
            // Notify the server about file changes to files in the workspace
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*')
        }
    };

    // Create the language client
    client = new LanguageClient(
        'exampleLsp',
        'Example LSP',
        serverOptions,
        clientOptions
    );

    // Set trace level
    client.trace = trace === 'verbose' ? 2 : trace === 'messages' ? 1 : 0;

    // Start the client (this will also launch the server)
    client.start();

    // Register commands
    context.subscriptions.push(
        vscode.commands.registerCommand('exampleLsp.restart', async () => {
            await client.stop();
            client.start();
            vscode.window.showInformationMessage('Example LSP server restarted');
        })
    );

    console.log('Example LSP client started');
}

/**
 * Deactivates the extension
 */
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
