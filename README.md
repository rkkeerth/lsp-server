# LSP Server

A Language Server Protocol (LSP) implementation providing language intelligence features for modern code editors and IDEs.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
- [Testing](#testing)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [License](#license)
- [Author](#author)

## Overview

This project implements a Language Server Protocol (LSP) server that provides language intelligence features such as code completion, go-to-definition, find references, and diagnostics for supported programming languages. The LSP enables a standardized communication protocol between code editors and language servers, allowing for rich editing experiences across different development environments.

### What is LSP?

The Language Server Protocol (LSP) is an open, JSON-RPC-based protocol that standardizes the communication between development tools and language servers. It was created by Microsoft and is now widely adopted across the developer ecosystem.

**Key Benefits:**
- Write language support once, use it in multiple editors
- Consistent developer experience across different IDEs
- Reduced maintenance overhead for language tooling

## Features

- **Code Completion**: Intelligent autocomplete suggestions based on context
- **Go to Definition**: Navigate to symbol definitions
- **Find References**: Locate all references to a symbol
- **Hover Information**: Display documentation and type information
- **Diagnostics**: Real-time error and warning detection
- **Code Actions**: Quick fixes and refactoring suggestions
- **Document Symbols**: Outline view of file structure
- **Workspace Symbols**: Project-wide symbol search
- **Formatting**: Code formatting support
- **Rename**: Symbol renaming with references update

## Prerequisites

Before installing and running the LSP server, ensure you have the following:

- **Runtime Environment**: [Specify required runtime, e.g., Node.js 16+, Python 3.8+, Go 1.19+]
- **Operating System**: Linux, macOS, or Windows
- **Memory**: Minimum 512MB RAM recommended
- **Supported Editors**: 
  - Visual Studio Code
  - Neovim
  - Emacs
  - Sublime Text
  - Any LSP-compatible editor

## Installation

### From Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/rkkeerth/lsp-server.git
   cd lsp-server
   ```

2. **Install dependencies:**
   ```bash
   # Example for Node.js-based server
   npm install
   
   # Example for Python-based server
   pip install -r requirements.txt
   
   # Example for Go-based server
   go mod download
   ```

3. **Build the server:**
   ```bash
   # Example build commands
   npm run build
   # or
   make build
   ```

### Using Package Manager

```bash
# Example for npm
npm install -g lsp-server

# Example for pip
pip install lsp-server
```

## Usage

### Starting the Server

#### Standalone Mode
```bash
# Start the LSP server on default port
lsp-server --stdio

# Start with custom configuration
lsp-server --stdio --config /path/to/config.json
```

#### Socket Mode
```bash
# Start server listening on a specific port
lsp-server --port 9257
```

### Editor Integration

#### Visual Studio Code

1. Install the companion extension or configure manually:

```json
{
  "languageServerExample.server": {
    "command": "lsp-server",
    "args": ["--stdio"]
  }
}
```

#### Neovim

Add to your `init.lua`:

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

configs.lsp_server = {
  default_config = {
    cmd = {'lsp-server', '--stdio'},
    filetypes = {'your-language'},
    root_dir = lspconfig.util.root_pattern('.git'),
  }
}

lspconfig.lsp_server.setup{}
```

#### Emacs

Add to your Emacs configuration:

```elisp
(use-package lsp-mode
  :hook ((your-mode . lsp))
  :commands lsp
  :config
  (lsp-register-client
   (make-lsp-client :new-connection (lsp-stdio-connection "lsp-server")
                    :major-modes '(your-mode)
                    :server-id 'lsp-server)))
```

### Basic Example

Once configured in your editor, the LSP server will automatically:

1. Analyze your code in real-time
2. Provide completion suggestions as you type
3. Display diagnostics for errors and warnings
4. Enable navigation features (go-to-definition, find references)

## Configuration

### Configuration File

Create a configuration file `lsp-config.json`:

```json
{
  "server": {
    "port": 9257,
    "host": "localhost",
    "logLevel": "info"
  },
  "language": {
    "diagnostics": {
      "enabled": true,
      "severity": "warning"
    },
    "completion": {
      "enabled": true,
      "triggerCharacters": [".", ":", ">"]
    },
    "formatting": {
      "enabled": true,
      "tabSize": 2,
      "insertSpaces": true
    }
  },
  "workspace": {
    "maxFiles": 10000,
    "excludePatterns": [
      "**/node_modules/**",
      "**/dist/**",
      "**/.git/**"
    ]
  }
}
```

### Environment Variables

- `LSP_SERVER_PORT`: Server port (default: 9257)
- `LSP_SERVER_LOG_LEVEL`: Logging level (debug, info, warn, error)
- `LSP_SERVER_CACHE_DIR`: Directory for caching parsed files
- `LSP_SERVER_MAX_WORKERS`: Maximum number of worker threads

### Command Line Options

```bash
lsp-server [options]

Options:
  --stdio                Use stdio for communication (default)
  --port <number>        Port number for socket communication
  --host <string>        Host address for socket communication
  --config <path>        Path to configuration file
  --log-file <path>      Path to log file
  --log-level <level>    Logging level (debug|info|warn|error)
  --version              Display version information
  --help                 Display help information
```

## Development

### Setting Up Development Environment

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/rkkeerth/lsp-server.git
   cd lsp-server
   ```

2. **Install development dependencies:**
   ```bash
   npm install --dev
   # or
   pip install -r requirements-dev.txt
   ```

3. **Run in development mode:**
   ```bash
   npm run dev
   # or
   python -m lsp_server --debug
   ```

### Project Structure

```
lsp-server/
├── src/                  # Source code
│   ├── server/          # LSP server implementation
│   ├── handlers/        # Request handlers
│   ├── parser/          # Language parser
│   ├── analyzer/        # Code analysis
│   └── utils/           # Utility functions
├── tests/               # Test files
├── docs/                # Documentation
├── examples/            # Usage examples
├── config/              # Configuration files
└── scripts/             # Build and utility scripts
```

### Code Style

This project follows standard coding conventions:

- Use consistent indentation (2 spaces)
- Follow language-specific style guides
- Write meaningful variable and function names
- Document public APIs and complex logic
- Keep functions focused and concise

### Building from Source

```bash
# Install dependencies
make deps

# Run linter
make lint

# Build the project
make build

# Build for specific platform
make build-linux
make build-macos
make build-windows
```

## Testing

### Running Tests

```bash
# Run all tests
npm test
# or
make test

# Run specific test suite
npm test -- --grep "completion"

# Run with coverage
npm run test:coverage

# Run integration tests
npm run test:integration
```

### Test Structure

```
tests/
├── unit/               # Unit tests
├── integration/        # Integration tests
├── e2e/               # End-to-end tests
└── fixtures/          # Test fixtures and data
```

### Writing Tests

Example test case:

```javascript
describe('Completion Handler', () => {
  it('should provide completion suggestions', async () => {
    const server = new LSPServer();
    const result = await server.handleCompletion({
      textDocument: { uri: 'file:///test.js' },
      position: { line: 10, character: 5 }
    });
    
    expect(result.items).to.have.length.greaterThan(0);
    expect(result.items[0]).to.have.property('label');
  });
});
```

## Architecture

### High-Level Architecture

```
┌─────────────────┐
│  Code Editor    │
│  (VSCode, Vim)  │
└────────┬────────┘
         │ LSP Protocol (JSON-RPC)
         │
┌────────▼────────┐
│   LSP Server    │
│                 │
│  ┌───────────┐  │
│  │  Handler  │  │
│  │  Manager  │  │
│  └─────┬─────┘  │
│        │        │
│  ┌─────▼─────┐  │
│  │  Parser   │  │
│  └─────┬─────┘  │
│        │        │
│  ┌─────▼─────┐  │
│  │ Analyzer  │  │
│  └─────┬─────┘  │
│        │        │
│  ┌─────▼─────┐  │
│  │  Cache    │  │
│  └───────────┘  │
└─────────────────┘
```

### Key Components

1. **Server Core**: Manages LSP lifecycle and communication
2. **Request Router**: Dispatches requests to appropriate handlers
3. **Document Manager**: Tracks open documents and changes
4. **Parser**: Parses source code into AST (Abstract Syntax Tree)
5. **Analyzer**: Performs semantic analysis and type checking
6. **Cache Manager**: Caches parsed results for performance
7. **Diagnostic Engine**: Generates errors and warnings

### Communication Flow

1. Editor sends LSP request (JSON-RPC)
2. Server validates and routes request
3. Handler processes request using parser/analyzer
4. Results are formatted per LSP specification
5. Response sent back to editor

## Contributing

We welcome contributions from the community! Here's how you can help:

### How to Contribute

1. **Fork the repository**
   ```bash
   # Click the 'Fork' button on GitHub
   git clone https://github.com/YOUR_USERNAME/lsp-server.git
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Write clean, documented code
   - Add tests for new features
   - Update documentation as needed

4. **Run tests and linting**
   ```bash
   make test
   make lint
   ```

5. **Commit your changes**
   ```bash
   git add .
   git commit -m "Add feature: your feature description"
   ```

6. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**
   - Go to the original repository on GitHub
   - Click 'New Pull Request'
   - Select your fork and branch
   - Provide a clear description of your changes

### Contribution Guidelines

- **Code Quality**: Follow existing code style and conventions
- **Testing**: Ensure all tests pass and add tests for new features
- **Documentation**: Update README and code comments as needed
- **Commit Messages**: Write clear, descriptive commit messages
- **Issue Tracking**: Reference related issues in your PR description
- **Small PRs**: Keep pull requests focused and reasonably sized

### Types of Contributions

- 🐛 Bug fixes
- ✨ New features
- 📝 Documentation improvements
- 🎨 Code style/formatting
- ⚡ Performance improvements
- ✅ Additional tests
- 🔧 Configuration improvements

### Reporting Issues

Found a bug or have a feature request?

1. Check if the issue already exists
2. Create a new issue with:
   - Clear, descriptive title
   - Detailed description
   - Steps to reproduce (for bugs)
   - Expected vs actual behavior
   - Environment details (OS, editor, versions)
   - Relevant logs or screenshots

### Code Review Process

- All submissions require review
- Maintainers will provide feedback
- Address review comments promptly
- Once approved, changes will be merged

## Troubleshooting

### Common Issues

#### Server Not Starting

**Problem**: Server fails to start or crashes immediately

**Solutions**:
- Check if the required runtime is installed and up to date
- Verify all dependencies are installed: `npm install` or equivalent
- Check for port conflicts if using socket mode
- Review server logs for specific error messages

#### No Completions or Features Working

**Problem**: LSP features not working in editor

**Solutions**:
- Ensure the server is running: check editor's LSP client status
- Verify correct file types are configured
- Check workspace root directory is correctly set
- Restart the LSP server and editor
- Enable debug logging to see communication between editor and server

#### High CPU or Memory Usage

**Problem**: Server consuming excessive resources

**Solutions**:
- Exclude large directories (node_modules, build outputs) in configuration
- Reduce `maxFiles` in workspace configuration
- Check for infinite loops in file watching
- Consider increasing cache size to reduce re-parsing

#### Diagnostics Not Showing

**Problem**: Errors and warnings not appearing

**Solutions**:
- Check if diagnostics are enabled in configuration
- Verify file is in workspace root
- Ensure file type is supported
- Check diagnostic severity settings

### Debug Mode

Enable verbose logging:

```bash
lsp-server --stdio --log-level debug --log-file /tmp/lsp-server.log
```

Then monitor the log file:
```bash
tail -f /tmp/lsp-server.log
```

### Getting Help

- 📖 Check the [documentation](docs/)
- 💬 Join our community discussions
- 🐛 Report bugs via [GitHub Issues](https://github.com/rkkeerth/lsp-server/issues)
- 📧 Contact the maintainers

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Keerth RK** ([@rkkeerth](https://github.com/rkkeerth))

---

**Project Links:**
- Repository: [https://github.com/rkkeerth/lsp-server](https://github.com/rkkeerth/lsp-server)
- Issue Tracker: [https://github.com/rkkeerth/lsp-server/issues](https://github.com/rkkeerth/lsp-server/issues)
- Documentation: [https://github.com/rkkeerth/lsp-server/docs](https://github.com/rkkeerth/lsp-server/docs)

---

*Built with ❤️ by the LSP Server team*
