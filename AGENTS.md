# AGENTS.md

This file provides guidance to AI coding agents (Gemini CLI, Codex, GPT-Engineer, etc.)
when working with code in this repository.

> [!IMPORTANT]
> **Whenever you change behaviour, add a feature, or modify a command's interface,
> you MUST update every document listed in the [Documentation Checklist](#documentation-checklist)
> that is affected by the change.** Do not leave docs stale.

---

## Build, Test, Run

```bash
make build        # build to bin/wazuh-cli (injects Version via ldflags from git describe)
make test         # go test -v ./...
make vet          # go vet ./...
make fmt          # go fmt ./...
make tidy         # go mod tidy
make lint         # golangci-lint run ./... (optional)

go test -v -run TestName ./internal/config   # single test
./bin/wazuh-cli --help                       # verify help text
```

`main.go` is a thin wrapper: `cmd.Execute()` then `os.Exit(cmd.HandleError(err))`.
All logic lives in `cmd/` and `internal/`.

---

## Architecture

### Execution flow
`cmd/root.go` defines a `PersistentPreRunE` that:
1. Parses `flagOverrides` into `config.Config`.
2. Resolves config in priority order: **CLI flags → `WAZUH_*` env → `.env` in CWD → `~/.config/wazuh/config.json`** (`internal/config/config.Load`).
3. Calls `cfg.Validate()` — requires URL + (user+password OR raw token).
4. Constructs `globalClient` (`*client.Client`) and `globalFmt` (`*output.Formatter`).

Subcommands read those three package-level globals (`globalCfg`, `globalClient`, `globalFmt`). No DI — add a new command by creating `cmd/<resource>.go`, registering in `init()` via `rootCmd.AddCommand(...)`, and using the globals.

**Pre-run skip list** (in `root.go`): commands that do not need an authenticated client must be excluded in `PersistentPreRunE` by name or parent name. Currently: `init`, `version`, `completion`, `config`, `auth`, `alert`. Extend this list for any new command that runs before a client exists (e.g. setup/auth/offline commands).

### Client (`internal/client/`)
- `client.go` — `http.Client` with TLS config (supports `--insecure`, custom CA, mTLS client cert/key).
- `auth.go` — `authManager` acquires a JWT via Basic Auth on `/security/user/authenticate`, caches it at `~/.config/wazuh/token` (0600), auto-refreshes on 401. Accepts a raw token via `--token`/`WAZUH_TOKEN` to bypass login.
- `request.go` / `response.go` — wrap HTTP calls; responses are typed errors implementing `IsAuth()` / `IsPermission()` / `IsAPI()` marker interfaces.

### Indexer (`internal/indexer/`)
Separate OpenSearch client used by `cmd/alert.go` to query Wazuh alerts. Config: `indexer_url`, `indexer_user`, `indexer_password`, `indexer_index` (or `WAZUH_INDEXER_*` env). If indexer user/password are unset, falls back to the manager's `user`/`password`.

### Output (`internal/output/`)
`Formatter.Write(data)` dispatches on format: `json` (default, pretty by default), `markdown` (renders a table via `table.go`), or `raw` (passthrough for `[]byte`/`json.RawMessage`). `ExitCode(err)` maps the marker interfaces above to exit codes **0/1/2/3/4** (ok / client / API / auth / permission). `WriteError` emits JSON envelopes on stdout in JSON mode, plain text on stderr otherwise.

---

## Key Patterns

- **Error contract**: JSON envelope `{"error": true, "code": N, "message": "...", "detail": "..."}`. To surface a specific exit code, return an error implementing the matching marker interface from `internal/output`.
- **0600 required**: config file and `~/.config/wazuh/config.json`'s sibling `token` must be 0600; loader warns on looser modes. `Config.Save` creates the dir as 0700 and writes as 0600.
- **Password from stdin**: use `config set password -P` (reads from stdin) or `config set password -p <value>` (inline flag). The global `--password -` flag on the root command also reads one line from stdin in `PersistentPreRunE`.
- **Setup**: prefer `wazuh-cli auth login` (interactive, no shell-history leak) over `wazuh-cli init` when helping users configure auth.
- **Sensitive keys**: `password` and `indexer_password` are masked in all display output; `config set` supports `-P`/`--from-stdin` and `-p`/`--password` flags to avoid leaking secrets.

---

## Project Layout

```
main.go
cmd/                Cobra commands; one file per API resource (agent, group, rule, decoder,
                    sca, vulnerability, syscollector, syscheck, rootcheck, active_response,
                    cluster, manager, security, cdb_list, task, mitre, logtest, alert,
                    auth, config, init, version). root.go owns persistent flags + globals.
internal/config/    4-tier config resolution, Validate, Save (0600).
internal/client/    HTTP client, JWT auth manager, typed request/response errors.
internal/indexer/   OpenSearch client for alert queries.
internal/output/    json/markdown/raw formatters, ExitCode mapping, table renderer.
docs/               User-facing docs (architecture, auth, configuration, commands, ...).
skill/SKILL.md      AI-agent skill definition for Claude Code / Gemini / Cline.
AGENTS.md           This file — guidance for AI coding agents.
CLAUDE.md           Guidance for Claude Code specifically.
```

---

## Documentation Checklist

> [!IMPORTANT]
> **Always update every applicable document below when making a change.** Stale docs
> mislead users and other agents. A pull request that changes behaviour without
> updating docs will be rejected.

| Document | Update when… |
|---|---|
| [`README.md`](README.md) | Any user-facing feature, flag, or workflow changes |
| [`CLAUDE.md`](CLAUDE.md) | Architecture changes, new patterns, pre-run skip list changes |
| [`AGENTS.md`](AGENTS.md) | Architecture changes, new patterns, any guidance for agents |
| [`skill/SKILL.md`](skill/SKILL.md) | Any command added/removed/changed, new security notes, new workflows |
| [`docs/configuration.md`](docs/configuration.md) | New config keys, changed `config set` behaviour, new flags |
| [`docs/authentication.md`](docs/authentication.md) | Auth flow changes, new credential input methods, security features |
| [`docs/commands.md`](docs/commands.md) | New subcommands or changed command signatures |
| [`docs/architecture.md`](docs/architecture.md) | Internal design changes, new packages, execution flow changes |
| [`docs/implementation.md`](docs/implementation.md) | Developer-facing implementation detail changes |
| [`docs/user-management.md`](docs/user-management.md) | RBAC or user management changes |

### Quick rule of thumb

- Changed a **flag or command interface**? → `README.md`, `skill/SKILL.md`, `docs/commands.md`, `docs/configuration.md`
- Changed **security or credential handling**? → `README.md`, `skill/SKILL.md`, `docs/authentication.md`, `CLAUDE.md`, `AGENTS.md`
- Changed **internal architecture**? → `CLAUDE.md`, `AGENTS.md`, `docs/architecture.md`, `docs/implementation.md`
- Added a **new command**? → All of the above, plus add an entry in `docs/commands.md`

---

## Security Rules (Never Violate)

1. **Never print or log raw passwords** — mask with `********` or `(already set)`.
2. **0600 file permissions** — enforce on both `config.json` and the token cache.
3. **Prefer stdin (`-P`) over inline flags (`-p`) for secrets** — document this order in help text.
4. **`auth login` is the gold standard** — recommend it first in any setup guidance.
5. **No interactive prompts in normal operation** — the tool must be fully scriptable by AI agents.
