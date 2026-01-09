# LSP Server

A comprehensive Language Server Protocol (LSP) implementation providing intelligent language features for modern code editors and IDEs.

## Table of Contents

- [About](#about)
- [What is LSP?](#what-is-lsp)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## About

This LSP server provides rich language intelligence features to enhance the development experience across multiple editors and IDEs. By implementing the Language Server Protocol, it enables consistent, powerful code editing capabilities without requiring custom integrations for each editor.

## What is LSP?

The **Language Server Protocol (LSP)** is an open, JSON-RPC-based protocol that standardizes communication between code editors (language clients) and language intelligence tools (language servers). 

### Why LSP?

Before LSP, each editor had to implement language support individually, leading to:
- Duplicated effort across editor teams
- Inconsistent feature availability
- Limited resources for niche programming languages

LSP solves this by:
- **Decoupling**: Separating language intelligence from editors
- **Reusability**: One server works with all LSP-compatible editors
- **Standardization**: Consistent features across different development environments

### Supported Editors

Any editor that supports LSP can use this server, including:
- Visual Studio Code
- Vim/Neovim
- Emacs
- Sublime Text
- Atom
- Eclipse
- IntelliJ IDEA
- And many more...

## Features

This LSP server provides the following language intelligence capabilities:

### Core Features

- **Code Completion (IntelliSense)**: Context-aware suggestions as you type
- **Go to Definition**: Navigate directly to where symbols are defined
- **Find References**: Locate all usages of a symbol across your codebase
- **Hover Information**: View documentation and type information on hover
- **Signature Help**: Display function/method signatures while typing parameters

### Diagnostics & Analysis

- **Error Detection**: Real-time syntax and semantic error checking
- **Warning Markers**: Highlight potential issues and code smells
- **Code Linting**: Enforce code quality and style guidelines

### Refactoring & Navigation

- **Rename Symbol**: Safely rename variables, functions, and classes across files
- **Document Symbols**: Outline view of all symbols in the current file
- **Workspace Symbols**: Search for symbols across the entire project
- **Code Actions**: Quick fixes and automated refactorings

### Formatting

- **Document Formatting**: Format entire files according to style rules
- **Range Formatting**: Format selected code sections
- **Format on Type**: Automatic formatting as you type

## Installation

### Prerequisites

Before installing, ensure you have the following:

```bash
# Example prerequisites (adjust based on your implementation)
- Node.js >= 14.x (or your runtime requirement)
- npm >= 6.x or yarn >= 1.22
```

### From Package Manager

```bash
# npm
npm install -g lsp-server

# yarn
yarn global add lsp-server

# pip (for Python-based implementations)
pip install lsp-server
```

### From Source

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Install dependencies
npm install  # or yarn install, pip install -r requirements.txt, etc.

# Build (if necessary)
npm run build

# Link globally for development
npm link
```

### Binary Installation

Download pre-built binaries from the [releases page](https://github.com/rkkeerth/lsp-server/releases) and add them to your PATH.

## Usage

### Getting Started

#### VS Code

Create or edit your `.vscode/settings.json`:

```json
{
  "lsp-server.enable": true,
  "lsp-server.trace.server": "verbose"
}
```

#### Neovim

Add to your Neovim configuration:

```lua
require'lspconfig'.lsp_server.setup{
  cmd = {"lsp-server", "--stdio"},
  filetypes = {"your-language"},
  root_dir = require'lspconfig'.util.root_pattern(".git", "package.json"),
}
```

#### Vim (with vim-lsp)

```vim
if executable('lsp-server')
  au User lsp_setup call lsp#register_server({
    \ 'name': 'lsp-server',
    \ 'cmd': {server_info->['lsp-server', '--stdio']},
    \ 'allowlist': ['your-language'],
    \ })
endif
```

#### Emacs (with lsp-mode)

```elisp
(use-package lsp-mode
  :hook (your-language-mode . lsp)
  :commands lsp
  :config
  (lsp-register-client
   (make-lsp-client :new-connection (lsp-stdio-connection '("lsp-server" "--stdio"))
                    :major-modes '(your-language-mode)
                    :server-id 'lsp-server)))
```

### Command Line Options

```bash
# Start in stdio mode (default for most editors)
lsp-server --stdio

# Start with TCP socket
lsp-server --port=6009

# Enable debug logging
lsp-server --stdio --log-level=debug

# Show version
lsp-server --version

# Show help
lsp-server --help
```

### Basic Example

Once configured, the LSP server works automatically in the background:

1. **Open a file** in your supported language
2. **Start typing** - get instant code completion
3. **Hover** over symbols to see documentation
4. **Use editor shortcuts** for go-to-definition, find-references, etc.

## Configuration

### Configuration File

Create a configuration file in your project root:

**`.lsp-server.json`**

```json
{
  "diagnostics": {
    "enabled": true,
    "level": "warning",
    "debounceTime": 500
  },
  "completion": {
    "enabled": true,
    "triggerCharacters": [".", ":", ">"],
    "maxItems": 100
  },
  "formatting": {
    "enabled": true,
    "tabSize": 2,
    "insertSpaces": true
  },
  "hover": {
    "enabled": true,
    "includeDocumentation": true
  },
  "workspace": {
    "maxNumberOfProblems": 1000,
    "excludePatterns": ["**/node_modules/**", "**/dist/**"]
  }
}
```

### Environment Variables

Configure the server using environment variables:

```bash
# Set log level
export LSP_SERVER_LOG_LEVEL=debug

# Set maximum memory usage
export LSP_SERVER_MAX_MEMORY=4096

# Custom configuration path
export LSP_SERVER_CONFIG_PATH=/path/to/config.json

# Enable performance profiling
export LSP_SERVER_PROFILE=true
```

### Editor-Specific Settings

#### VS Code

```json
{
  "lsp-server.enable": true,
  "lsp-server.diagnostics.enable": true,
  "lsp-server.trace.server": "messages",
  "lsp-server.completion.autoImport": true,
  "lsp-server.formatting.provider": "lsp-server"
}
```

#### Neovim/Vim

```lua
require'lspconfig'.lsp_server.setup{
  settings = {
    lspServer = {
      diagnostics = { enabled = true },
      completion = { autoImport = true },
      formatting = { enabled = true }
    }
  }
}
```

## Development

### Setting Up the Development Environment

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Install development dependencies
npm install

# Run in development mode with auto-reload
npm run dev

# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Lint code
npm run lint

# Format code
npm run format
```

### Project Structure

```
lsp-server/
├── src/                  # Source code
│   ├── server.ts        # Main server implementation
│   ├── handlers/        # LSP request handlers
│   ├── analyzer/        # Code analysis engine
│   ├── completion/      # Completion provider
│   ├── diagnostics/     # Diagnostic engine
│   └── utils/           # Utility functions
├── tests/               # Test files
│   ├── unit/           # Unit tests
│   ├── integration/    # Integration tests
│   └── fixtures/       # Test fixtures
├── docs/               # Documentation
├── examples/           # Example configurations
├── package.json        # Package metadata
├── tsconfig.json       # TypeScript configuration
└── README.md          # This file
```

### Building

```bash
# Build for production
npm run build

# Build with watch mode
npm run build:watch

# Clean build artifacts
npm run clean
```

### Testing

```bash
# Run all tests
npm test

# Run specific test file
npm test -- --testPathPattern=completion

# Run tests in watch mode
npm run test:watch

# Run integration tests
npm run test:integration

# Generate coverage report
npm run test:coverage
```

### Debugging

#### VS Code Debug Configuration

Add to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "node",
      "request": "launch",
      "name": "Debug LSP Server",
      "program": "${workspaceFolder}/src/server.ts",
      "args": ["--stdio"],
      "outFiles": ["${workspaceFolder}/dist/**/*.js"],
      "sourceMaps": true,
      "console": "integratedTerminal"
    }
  ]
}
```

#### Logging

Enable verbose logging:

```bash
LSP_SERVER_LOG_LEVEL=debug lsp-server --stdio 2> lsp-server.log
```

### Code Quality

We maintain high code quality standards:

```bash
# Run linter
npm run lint

# Fix linting issues automatically
npm run lint:fix

# Check TypeScript types
npm run type-check

# Run all quality checks
npm run check
```

## Contributing

We welcome contributions from the community! Here's how you can help:

### How to Contribute

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/lsp-server.git
   cd lsp-server
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** and commit them:
   ```bash
   git commit -m "Add feature: your feature description"
   ```
5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
6. **Open a Pull Request** on GitHub

### Contribution Guidelines

#### Code Style

- Follow the existing code style and conventions
- Run `npm run lint` before committing
- Write meaningful commit messages following [Conventional Commits](https://www.conventionalcommits.org/)
- Keep functions small and focused
- Add comments for complex logic

#### Testing

- Write tests for new features
- Ensure all tests pass before submitting PR
- Maintain or improve code coverage
- Include both unit and integration tests where applicable

#### Documentation

- Update README.md if adding new features
- Add inline code documentation for public APIs
- Include examples for new functionality
- Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/)

#### Pull Request Process

1. Update documentation to reflect changes
2. Add tests that prove your fix/feature works
3. Ensure CI/CD pipeline passes
4. Request review from maintainers
5. Address review feedback promptly
6. Squash commits if requested

### Reporting Bugs

Found a bug? Please [open an issue](https://github.com/rkkeerth/lsp-server/issues/new) with:

- **Clear title** describing the issue
- **Steps to reproduce** the problem
- **Expected behavior** vs actual behavior
- **Environment details** (OS, editor, versions)
- **Logs or error messages** if available
- **Minimal reproduction** example if possible

### Feature Requests

Have an idea? [Open a feature request](https://github.com/rkkeerth/lsp-server/issues/new) with:

- **Clear description** of the feature
- **Use case** explaining why it's needed
- **Proposed solution** if you have one
- **Alternative approaches** you've considered

### Code of Conduct

This project follows a Code of Conduct to ensure a welcoming environment:

- **Be respectful** of different viewpoints and experiences
- **Accept constructive criticism** gracefully
- **Focus on what is best** for the community
- **Show empathy** towards other community members

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### MIT License Summary

```
Copyright (c) 2026 rkkeerth

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## Additional Resources

### Learning Resources

- [Official LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [LSP Documentation](https://microsoft.github.io/language-server-protocol/specifications/specification-current/)
- [Language Server Protocol Wiki](https://en.wikipedia.org/wiki/Language_Server_Protocol)

### Community

- **Issues**: [GitHub Issues](https://github.com/rkkeerth/lsp-server/issues)
- **Discussions**: [GitHub Discussions](https://github.com/rkkeerth/lsp-server/discussions)
- **Twitter**: [@rkkeerth](https://twitter.com/rkkeerth)

### Related Projects

- [Official LSP Libraries](https://microsoft.github.io/language-server-protocol/implementors/sdks/)
- [LSP Example Implementations](https://github.com/microsoft/vscode-extension-samples/tree/main/lsp-sample)

---

**Made with ❤️ by [rkkeerth](https://github.com/rkkeerth)**

*Star this repository if you find it useful! ⭐*
