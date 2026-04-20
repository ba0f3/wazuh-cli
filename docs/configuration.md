# Configuration System

`wazuh-cli` uses a 4-tier configuration resolution system. Settings are merged in the following priority (highest to lowest).

## Resolution Priority

1. **CLI Flags**: (e.g., `--url`, `--user`) Always take precedence.
2. **Environment Variables**: Prefixed with `WAZUH_` (e.g., `WAZUH_URL`, `WAZUH_PASSWORD`).
3. **Dotenv File**: A `.env` file in the current working directory.
4. **Config File**: Local JSON file at `~/.config/wazuh/config.json`.

## Configuration Keys

| Key | Env Var | Description |
|---|---|---|
| `url` | `WAZUH_URL` | Base URL of the Wazuh API (e.g., `https://localhost:55000`) |
| `user` | `WAZUH_USER` | API Username |
| `password` | `WAZUH_PASSWORD` | API Password |
| `insecure` | `WAZUH_INSECURE` | Skip TLS verification (boolean: `true`/`false`) |
| `output` | `WAZUH_OUTPUT` | Default output format (`json`, `markdown`, `raw`) |
| `timeout` | `WAZUH_TIMEOUT` | HTTP timeout in seconds |
| `debug` | `WAZUH_DEBUG` | Enable verbose debug logging |
| `indexer_url` | `WAZUH_INDEXER_URL` | Wazuh Indexer URL (e.g. `https://indexer:9200`) |
| `indexer_user` | `WAZUH_INDEXER_USER` | Indexer username (defaults to `user` if omitted) |
| `indexer_password` | `WAZUH_INDEXER_PASSWORD` | Indexer password (defaults to `password` if omitted) |
| `indexer_index` | `WAZUH_INDEXER_INDEX` | Index pattern for alerts (defaults to `wazuh-alerts-4.x-*`) |

> **Note on Indexer Credentials**: If you configure `indexer_url` to query alerts but leave `indexer_user` and `indexer_password` blank, `wazuh-cli` assumes your Indexer and Manager API share the same credentials and will use `user` and `password` for Indexer authentication.

## Config Management Commands

### `wazuh-cli config list`
Shows all values currently saved in the JSON configuration file. Secrets are masked.

### `wazuh-cli config set <KEY> [VALUE]`
Updates a specific key in the configuration file.

```bash
wazuh-cli config set insecure true
wazuh-cli config set url https://wazuh:55000
```

For the sensitive keys `password` and `indexer_password`, two explicit flags
control how the secret is supplied (to keep it out of shell history):

| Flag | Short | Behaviour |
|---|---|---|
| `--from-stdin` | `-P` | Read value from stdin — safest, never in history or process list |
| `--password` | `-p` | Pass value inline as a flag argument |

If neither flag is given, the positional argument is used as a fallback
(works for all keys, but discouraged for passwords).

```bash
# Recommended: read from stdin
wazuh-cli config set password -P < /run/secrets/wazuh_pass
read -rs PASS && printf '%s' "$PASS" | wazuh-cli config set password -P

# Inline flag (not in shell history, but visible in process list)
wazuh-cli config set password -p s3cr3t

# Works for indexer_password too
wazuh-cli config set indexer_password -P < /run/secrets/indexer_pass
```

### `wazuh-cli config get <KEY>`
Retrieves a specific value. If the value is empty, it reports "not set".

### `wazuh-cli config delete <KEY>`
Resets a value to its default (empty/false).

## Location
The default configuration directory is `~/.config/wazuh/`. 
- Configuration: `config.json`
- Token Cache: `token`
- Custom Path: You can override the config path using the `--config` flag.
