"""Entry point for the LSP server."""
import sys

from lsp_server.server import Server


def main() -> None:
    server = Server(sys.stdin.buffer, sys.stdout.buffer)
    try:
        server.run()
    except SystemExit:
        raise
    except Exception:
        sys.exit(1)


if __name__ == "__main__":
    main()
