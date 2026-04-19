# Implementation Details

This document describes the technical implementation of the `wazuh-cli` codebase, including package structure, patterns used, and error handling strategies.

## Package Structure

The project follows the standard Go layout:

- `cmd/`: Command definitions using the [Cobra](https://github.com/spf13/cobra) framework.
  - `root.go`: Defines the base command, global flags, and the `PersistentPreRunE` hook for configuration resolution and client initialization.
  - Resource files (e.g., `agent.go`, `rule.go`): Subcommands for specific Wazuh API resources.
- `internal/config/`: Configuration management.
  - Uses `json` for persistence in `~/.config/wazuh/config.json`.
  - Implements a hierarchical resolver that merges CLI flags, environment variables, `.env` files, and the JSON config.
- `internal/client/`: The Wazuh API client.
  - `authManager`: Handles JWT acquisition, multi-method fallbacks (GET/POST/JSON), and encrypted-on-disk token caching.
  - `Client`: Wraps `http.Client` with custom TLS configurations (custom CA, mTLS, insecure override).
- `internal/output/`: Standardized output engine.
  - Handles conversion between API responses and JSON, Markdown, or Raw formats.
  - Maps API errors to standardized exit codes.

## Authentication Logic

The authentication system (`internal/client/auth.go`) is designed to be resilient across different Wazuh versions:

1.  **GET Basic Auth**: Attempted first as per modern Wazuh API specs.
2.  **POST Basic Auth**: Attempted if GET fails (compatibility with some proxies).
3.  **POST JSON Body**: Final fallback for stricter API implementations.

Tokens are cached locally in `~/.config/wazuh/token` to minimize re-authentication overhead.

## Error Handling

Errors are treated as first-class citizens:

- **Exit Codes**:
  - `0`: Success
  - `1`: Client/Validation error
  - `2`: API error (Wazuh returned an error envelope)
  - `3`: Authentication failure
  - `4`: RBAC / Permission denied
- **JSON Error Envelopes**: When running with `--output json`, errors are returned as machine-readable JSON objects to stderr, allowing AI agents and scripts to parse the failure reason.

## Output Formatting

The `internal/output` package uses Go templates and reflection to handle dynamic API responses:
- **JSON**: Uses `json.MarshalIndent` for human readability when `pretty` is enabled.
- **Markdown**: Automatically flattens nested JSON objects into a table format, making it ideal for human consumption and LLM report generation.

## Dependencies

The project maintains a minimal dependency footprint:
- `github.com/spf13/cobra`: CLI framework.
- `github.com/joho/godotenv`: `.env` file support.
- `golang.org/x/term`: Secure password input.
- `github.com/spf13/pflag`: Flag handling (dependency of Cobra).
