import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  // Get server path from configuration
  const config = workspace.getConfiguration('basicLspServer');
  const serverPath = config.get<string>('serverPath', 'lsp-server');

  // Define server options
  const serverOptions: ServerOptions = {
    command: serverPath,
    args: [],
    options: {}
  };

  // Define client options
  const clientOptions: LanguageClientOptions = {
    // Register the server for plain text documents
    documentSelector: [{ scheme: 'file', language: 'plaintext' }],
    synchronize: {
      // Notify the server about file changes to '.clientrc' files contained in the workspace
      fileEvents: workspace.createFileSystemWatcher('**/.clientrc')
    }
  };

  // Create the language client
  client = new LanguageClient(
    'basicLspServer',
    'Basic LSP Server',
    serverOptions,
    clientOptions
  );

  // Start the client (which also launches the server)
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
