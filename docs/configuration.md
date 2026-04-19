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

## Config Management Commands

### `wazuh-cli config list`
Shows all values currently saved in the JSON configuration file. Secrets are masked.

### `wazuh-cli config set <KEY> <VALUE>`
Updates a specific key in the configuration file.
```bash
wazuh-cli config set insecure true
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
