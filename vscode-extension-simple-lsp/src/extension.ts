// The module 'vscode' contains the VS Code extensibility API
//
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export async function activate(context: vscode.ExtensionContext) {
  const serverOpts: ServerOptions = {
    command: "/Users/douglasmendes/Git/personal/golang-lsp/main",
    args: [],
  };
  const clientOpts: LanguageClientOptions = {
    documentSelector: [
      {
        scheme: "file",
        language: "markdown",
      },
      {
        scheme: "file",
        language: "c",
      },
    ],
    synchronize: {
      fileEvents: vscode.workspace.createFileSystemWatcher('**/*'),
    },
  };

  // Use the console to output diagnostic information (console.log) and errors (console.error)
  // This line of code will only be executed once when your extension is activated
  console.log(
    'Congratulations, your extension "vscode-extension-simple-lsp" is now active!',
  );

  // The command has been defined in the package.json file
  // Now provide the implementation of the command with registerCommand
  // The commandId parameter must match the command field in package.json
  const disposable = vscode.commands.registerCommand(
    "vscode-extension-simple-lsp.helloWorld",
    () => {
      // The code you place here will be executed every time your command is executed
      // Display a message box to the user
      vscode.window.showInformationMessage(
        "Hello World from vscode-extension-simple-lsp!",
      );
    },
  );

  context.subscriptions.push(disposable);

  const client = new LanguageClient("simple-lsp", serverOpts, clientOpts, true);
  await client.start();
}

// This method is called when your extension is deactivated
export function deactivate() {}
