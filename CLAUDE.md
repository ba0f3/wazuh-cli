# wazuh-cli — Agent Context (CLAUDE.md)

This repository contains `wazuh-cli`, a Go CLI for the Wazuh Server API.

## Quick Start

```bash
# Setup
wazuh-cli init --url https://wazuh:55000 --user admin --password secret

# Verify
wazuh-cli manager info

# List active agents
wazuh-cli agent list --status active

# Get specific agent
wazuh-cli agent get 001
```

## Config Location

`~/.config/wazuh/config.json` — created by `wazuh-cli init`.

## Build & Run

```bash
make build        # build to bin/wazuh-cli
make test         # run tests
make fmt          # format code
make tidy         # go mod tidy
go build -o wazuh-cli .
./wazuh-cli --help
```

## Project Layout

```
main.go             → entry point
cmd/                → Cobra commands (one file per resource)
  root.go           → persistent flags, config loading, client init
  agent.go          → agent management
  group.go          → agent groups
  rule.go           → rules
  decoder.go        → decoders
  sca.go            → Security Config Assessment
  vulnerability.go  → vulnerability detection
  syscollector.go   → system inventory
  syscheck.go       → FIM/syscheck
  rootcheck.go      → rootcheck
  active_response.go → active response
  cluster.go        → cluster management
  manager.go        → server info/control
  security.go       → RBAC (users, roles, policies)
  cdb_list.go       → CDB lists
  task.go           → tasks
  mitre.go          → MITRE ATT&CK
  logtest.go        → log testing
  init.go           → setup wizard

internal/
  config/config.go  → 4-tier config resolution (flags>env>.env>file)
  client/           → HTTP client, JWT auth, request/response
  output/           → JSON/Markdown/raw formatters

skill/SKILL.md      → AI agent skill documentation
```

## Key Patterns

- All commands return JSON by default (`--output markdown` for tables)
- Errors are JSON envelopes: `{"error": true, "code": N, "message": "..."}`
- Exit codes: 0=ok, 1=client err, 2=API err, 3=auth fail, 4=permission denied
- Auth: JWT acquired via Basic Auth, cached in `~/.config/wazuh/token`
- Config priority: CLI flags > WAZUH_* env vars > .env > ~/.config/wazuh/config.json

## Testing Changes

```bash
go build ./...    # must succeed
go vet ./...      # must pass
./wazuh-cli --help  # verify help text
```
