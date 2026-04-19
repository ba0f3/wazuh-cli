# wazuh-cli

> AI-agent-first CLI for the Wazuh Server API

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

`wazuh-cli` wraps the entire Wazuh Server REST API in a single binary with structured JSON/Markdown output, making it the ideal tool for AI agents (Claude Code, Gemini CLI, Cline, etc.) and human operators alike.

---

## Features

- **Complete API coverage** — agents, rules, decoders, SCA, vulnerabilities, FIM, syscollector, active response, cluster, RBAC, MITRE ATT&CK, and more
- **AI Agent-First**: Machine-parsable JSON output, deterministic exit codes, and a built-in Agent Skill.
- **Secure**: 0600 file permissions, credential masking, and secure interactive login.

## Documentation

- [Architecture](docs/architecture.md) - Design philosophy and component structure
- [Authentication](docs/authentication.md) - Multi-method auth, token caching, and security
- [Configuration](docs/configuration.md) - Priority resolution and config management
- [Command Reference](docs/commands.md) - Full list of supported API resources
- [Implementation](docs/implementation.md) - Technical details and package structure

## Installation

```bash
# From source
git clone https://github.com/ba0f3/wazuh-cli
cd wazuh-cli
go build -o wazuh-cli .
sudo mv wazuh-cli /usr/local/bin/
```

---

## Quick Start

```bash
# 1. Initialize configuration
wazuh-cli init --url https://wazuh-server:55000 --user admin --password secret

# 2. Verify connectivity
wazuh-cli manager info

# 3. List active agents
wazuh-cli agent list --status active

# 4. Find critical vulnerabilities
wazuh-cli vulnerability list 001 --severity critical

# 5. Check cluster health
wazuh-cli cluster health
```

---

## Configuration

Configuration is loaded in priority order (highest wins):

| Priority | Source | Example |
|---|---|---|
| 1 | CLI flag | `--url https://wazuh:55000` |
| 2 | Environment variable | `WAZUH_URL=https://wazuh:55000` |
| 3 | `.env` in current directory | `WAZUH_URL=https://wazuh:55000` |
| 4 | Config file | `~/.config/wazuh/config.json` |

### Config File (`~/.config/wazuh/config.json`)

```json
{
  "url": "https://wazuh-server:55000",
  "user": "admin",
  "password": "your-password",
  "insecure": false,
  "timeout": 30,
  "output": "json",
  "pretty": true
}
```

> **Security**: Config file must have `0600` permissions (enforced on load).

### Environment Variables

| Variable | Description |
|---|---|
| `WAZUH_URL` | API URL |
| `WAZUH_USER` | Username |
| `WAZUH_PASSWORD` | Password |
| `WAZUH_TOKEN` | Raw JWT (skips user/password auth) |
| `WAZUH_INSECURE` | Skip TLS (`true`/`1`) |
| `WAZUH_CA_CERT` | Path to custom CA cert |

---

## Command Reference

```
wazuh-cli
├── init                    Setup configuration
├── auth                    Authentication management
│   └── login               Secure interactive login
├── config                  Manage local configuration file
│   ├── list                List all configured values
│   ├── get <KEY>           Get a specific value
│   ├── set <KEY> <VALUE>   Set a value
│   └── delete <KEY>        Delete a value
├── agent                   Agent management
│   ├── list                List agents (--status, --group, --query)
│   ├── get <ID>            Get agent details
│   ├── delete <ID>         Delete agent
│   ├── restart <ID>        Restart agent
│   ├── summary             Agent status counts
│   ├── key <ID>            Get agent key
│   ├── upgrade <ID>        Upgrade agent
│   └── config <ID>         Get active config (--component, --configuration)
├── group                   Agent group management
│   ├── list / get / create / delete
│   ├── agents <NAME>       List group agents
│   ├── add-agent           Add agent to group
│   └── remove-agent        Remove agent from group
├── rule                    Rules
│   ├── list / get / files / groups
├── decoder                 Decoders
│   ├── list / get / files
├── sca                     Security Configuration Assessment
│   ├── list / get / checks
├── vulnerability           Vulnerability detection
│   ├── list / summary
├── syscollector            System inventory
│   ├── hardware / os / packages / processes / ports / netaddr / netiface / hotfixes
├── syscheck                File Integrity Monitoring
│   ├── list / run / clear / last-scan
├── rootcheck               Rootcheck
│   ├── list / run / clear / last-scan
├── active-response         Active response
│   └── run
├── cluster                 Cluster management
│   ├── status / nodes / node / health / config / restart
├── manager                 Manager info & control
│   ├── info / status / config / stats / logs / log-summary / restart / validation
├── security                RBAC
│   ├── user  (list/get/create/update/delete)
│   ├── role  (list/get/create/delete)
│   ├── policy (list/get/delete)
│   └── rule  (list/get)
├── cdb-list                CDB lists
│   ├── list / get / delete
├── task                    Tasks
│   └── list
├── mitre                   MITRE ATT&CK
│   ├── list / get
├── logtest                 Log test engine
│   ├── run / end
└── ciscat                  CIS-CAT results
    └── results <AGENT_ID>
```

---

## Global Flags

```
-u, --url string         Wazuh API URL
    --user string        API username
    --password string    API password (use '-' to read from stdin)
    --token string       Raw JWT token
-o, --output string      Output format: json, markdown, raw (default: json)
    --pretty             Pretty-print JSON
-k, --insecure           Skip TLS certificate verification
    --ca-cert string     Custom CA certificate path
    --timeout int        HTTP timeout in seconds (default: 30)
    --debug              Debug logging to stderr
-q, --quiet              Suppress informational messages
    --config string      Config file path
```

---

## AI Agent Integration

Install the built-in skill for Claude Code or compatible agents:

```bash
# For Claude Code (copy to your skills directory)
cp skill/SKILL.md ~/.claude/skills/wazuh-cli.md

# Or reference the skill directly in your CLAUDE.md:
# See skill/SKILL.md for wazuh-cli usage
```

The `skill/SKILL.md` file contains:
- Complete command reference with examples
- Common investigation workflows
- Error handling patterns
- Output format guidance for programmatic parsing

---

## Security

- Config file permissions enforced to `0600`
- JWT tokens cached at `~/.config/wazuh/token` (`0600`)
- Credentials never logged (even in `--debug` mode)
- TLS verification on by default (`--insecure` explicitly required)
- Password from stdin: `--password -`

---

## License

MIT
