# LSP Server - Setup & Verification Checklist

Use this checklist to verify your LSP server installation and setup.

## ✅ Initial Setup

### Prerequisites
- [ ] Go 1.21+ installed (`go version`)
- [ ] Git installed (`git --version`)
- [ ] Terminal/command line access
- [ ] Text editor (optional: VS Code for extension testing)

### Repository Setup
- [ ] Repository cloned or downloaded
- [ ] Current directory is `lsp-server/`
- [ ] All files present (run `ls -la`)

## ✅ File Verification

### Core Go Files (7 files)
- [ ] `main.go` - Entry point
- [ ] `go.mod` - Module definition
- [ ] `internal/lsp/protocol.go` - Protocol types
- [ ] `internal/lsp/server.go` - Server implementation
- [ ] `internal/lsp/lifecycle.go` - Lifecycle handlers
- [ ] `internal/lsp/textdocument.go` - Document handlers
- [ ] `internal/lsp/server_test.go` - Unit tests

### Documentation Files (7+ files)
- [ ] `README.md` - Main documentation
- [ ] `QUICKSTART.md` - Quick start guide
- [ ] `ARCHITECTURE.md` - Architecture documentation
- [ ] `TESTING.md` - Testing guide
- [ ] `CONTRIBUTING.md` - Contribution guidelines
- [ ] `PROJECT_SUMMARY.md` - Project overview
- [ ] `LICENSE` - MIT License

### Build & Configuration Files
- [ ] `Makefile` - Build automation
- [ ] `.gitignore` - Git ignore patterns

### Example Files
- [ ] `examples/test_client.go` - Test client
- [ ] `examples/vscode-extension/package.json` - VS Code extension
- [ ] `examples/vscode-extension/src/extension.ts` - Extension code

## ✅ Build Verification

### Step 1: Check Go Environment
```bash
go version
# Should show: go version go1.21 or higher
```
- [ ] Go version is 1.21 or higher

### Step 2: Check Module
```bash
cat go.mod
# Should show: module github.com/rkkeerth/lsp-server
```
- [ ] Module path is correct

### Step 3: Download Dependencies (if any)
```bash
go mod tidy
```
- [ ] Command completes without errors
- [ ] No dependencies needed (this project uses standard library only)

### Step 4: Build the Server
```bash
make build
# OR
go build -o lsp-server .
```
- [ ] Build completes without errors
- [ ] `lsp-server` binary exists
- [ ] Binary is executable (`ls -lh lsp-server`)

## ✅ Testing Verification

### Unit Tests
```bash
make test
# OR
go test ./...
```
- [ ] All tests pass
- [ ] No test failures
- [ ] Test coverage reported

Expected output:
```
ok      github.com/rkkeerth/lsp-server/internal/lsp     0.123s
```

### Integration Test
```bash
cd examples
go run test_client.go ../lsp-server
```
- [ ] Test client starts
- [ ] Initialize test passes
- [ ] didOpen test passes
- [ ] didChange test passes
- [ ] didClose test passes
- [ ] Shutdown test passes
- [ ] "All tests completed successfully!" message appears

## ✅ Code Quality Checks

### Format Check
```bash
make fmt
# OR
gofmt -l .
```
- [ ] No files need formatting
- [ ] Output is empty or all files formatted

### Static Analysis
```bash
make vet
# OR
go vet ./...
```
- [ ] No issues reported
- [ ] Command completes successfully

## ✅ Feature Verification

### LSP Protocol Features
- [ ] JSON-RPC 2.0 message parsing works
- [ ] Content-Length header handling works
- [ ] Request/response cycle works
- [ ] Notification handling works

### Lifecycle Features
- [ ] Initialize request accepted
- [ ] Server capabilities returned
- [ ] Shutdown request handled
- [ ] Exit notification handled

### Document Synchronization Features
- [ ] textDocument/didOpen tracked
- [ ] textDocument/didChange processed
- [ ] textDocument/didClose removed
- [ ] Multiple documents supported

### Concurrency Features
- [ ] Multiple requests handled concurrently
- [ ] Thread-safe document access
- [ ] No race conditions (run: `go test -race ./...`)

## ✅ Documentation Verification

### README.md
- [ ] Features section complete
- [ ] Installation instructions clear
- [ ] Usage examples provided
- [ ] Architecture overview present

### QUICKSTART.md
- [ ] 5-minute setup guide present
- [ ] Test instructions clear
- [ ] Common issues addressed

### ARCHITECTURE.md
- [ ] Component descriptions present
- [ ] Design decisions explained
- [ ] Diagrams or visualizations included

### TESTING.md
- [ ] Multiple testing approaches documented
- [ ] Test scenarios provided
- [ ] Debugging tips included

## ✅ Optional: VS Code Extension Test

### Extension Setup
```bash
cd examples/vscode-extension
npm install
npm run compile
```
- [ ] Dependencies installed
- [ ] TypeScript compiles without errors
- [ ] `out/extension.js` exists

### Extension Testing
- [ ] Open `examples/vscode-extension` in VS Code
- [ ] Press F5 to launch Extension Development Host
- [ ] Create a `.txt` file in the new window
- [ ] Type some text
- [ ] Check Output panel → "Basic LSP Server"

Extension Verification:
- [ ] Extension activates on .txt files
- [ ] Server starts automatically
- [ ] didOpen notification sent
- [ ] didChange notifications sent while typing
- [ ] No errors in Output panel

## ✅ Manual Server Test

### Start Server
```bash
./lsp-server
```
- [ ] Server starts without errors
- [ ] Logs appear on stderr
- [ ] "LSP Server starting..." message appears

### Send Initialize (in another terminal)
```bash
# This is a manual test - requires careful message formatting
echo -e 'Content-Length: 123\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":null,"rootUri":"file:///tmp","capabilities":{}}}' | ./lsp-server
```
- [ ] Server responds with initialize result
- [ ] Capabilities are listed
- [ ] No errors in logs

## ✅ Performance Checks

### Memory Usage
```bash
# Start server in background
./lsp-server &
PID=$!
ps -p $PID -o rss,vsz,cmd
```
- [ ] RSS (resident memory) < 50 MB initially
- [ ] Memory usage reasonable

### Response Time
- [ ] Initialize responds < 100ms
- [ ] Document operations respond < 10ms
- [ ] No noticeable lag

## ✅ Error Handling Verification

### Invalid Messages
Test that server handles errors gracefully:
- [ ] Malformed JSON doesn't crash server
- [ ] Unknown methods return MethodNotFound error
- [ ] Invalid params return InvalidParams error
- [ ] Server continues after errors

### State Validation
- [ ] Can't initialize twice
- [ ] Can't send didChange before initialize
- [ ] Shutdown/exit sequence validated

## ✅ Final Checks

### Project Completeness
- [ ] All features implemented as specified
- [ ] All documentation complete
- [ ] All tests passing
- [ ] No critical TODOs remaining

### Code Quality
- [ ] Code formatted consistently
- [ ] No lint errors
- [ ] Test coverage > 70%
- [ ] All exported functions documented

### Repository State
- [ ] .gitignore properly configured
- [ ] No build artifacts in git
- [ ] README.md accurate and complete
- [ ] LICENSE file present

## ✅ Ready for Use!

If all checks pass:
- ✅ LSP server is fully functional
- ✅ Tests are passing
- ✅ Documentation is complete
- ✅ Ready for development or extension

## Troubleshooting

### Build Fails
**Issue**: `go: command not found`  
**Solution**: Install Go from https://golang.org/dl/

**Issue**: `package X is not in GOROOT`  
**Solution**: Run `go mod tidy`

### Tests Fail
**Issue**: Test failures  
**Solution**: Check error messages, ensure no file conflicts

**Issue**: Integration test can't find server  
**Solution**: Build server first with `make build`

### Server Won't Start
**Issue**: `Permission denied`  
**Solution**: `chmod +x lsp-server`

**Issue**: No output  
**Solution**: Check stderr with `./lsp-server 2>log.txt`

## Need Help?

- 📖 Check **README.md** for detailed documentation
- 🏗️ Review **ARCHITECTURE.md** for design details  
- 🧪 See **TESTING.md** for testing help
- 💬 Open an issue on GitHub
- 📧 Contact maintainers

## Success Criteria

✅ **Minimum Requirements:**
- Server builds successfully
- Unit tests pass
- Integration test passes
- Documentation accessible

✅ **Full Success:**
- All checklist items complete
- VS Code extension works
- Manual testing successful
- Performance acceptable

---

**Last Updated**: 2024  
**Version**: 0.1.0

Once all items are checked, your LSP server is ready to use! 🎉
