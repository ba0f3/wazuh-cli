# Authentication and Security

`wazuh-cli` implements a robust authentication system to ensure secure and reliable communication with the Wazuh API.

## Authentication Flow

The CLI uses JWT (JSON Web Tokens) for all API requests. The flow for obtaining and using a token is:

1. **Check Cache**: Look for a valid JWT in `~/.config/wazuh/token`.
2. **Re-authenticate**: If no valid token exists, perform a fresh login using credentials.
3. **Persist**: Save the new JWT to the disk cache with `0600` permissions.

## Authentication Methods (Fallbacks)

Wazuh servers across different versions (4.x, 5.x) and configurations (Proxies, WAFs) handle authentication differently. To maximize compatibility, the CLI tries the following methods in order:

1. **GET with Basic Auth**: Standard for most Wazuh 4.x+ installations.
2. **POST with Basic Auth**: Used when proxies or older API versions expect POST.
3. **POST with JSON Body**: Fallback for APIs that require credentials in the request body.

## Security Features

### Credential Masking
- **Debug Logs**: When running with `--debug`, the password is NEVER printed. Instead, the password length is shown to help with troubleshooting.
- **Config Output**: `wazuh-cli config list` masks the `password` value.
- **Set Feedback**: `wazuh-cli config set password` confirms the update without echoing the secret.

### Interactive Login
The `wazuh-cli auth login` command uses a secure terminal prompt (masked input) to read passwords. This prevents secrets from being stored in shell history files (e.g., `.bash_history` or `.fish_history`) or exposed via process lists.

### Token Lifecycle
- Tokens are cached for 14 minutes (Wazuh default expiry is 15 minutes).
- `wazuh-cli auth logout` (planned) or deleting the token file manually invalidates the session locally.

## Certificate Handling
- **--insecure (-k)**: Disables TLS verification for servers using self-signed certificates.
- **--ca-cert**: Allows providing a custom Root CA for private certificate authorities.
- **mTLS**: Supports `--client-cert` and `--client-key` for environments requiring mutual TLS.

## OpenSearch Indexer Authentication

When running `wazuh-cli alert` commands, the CLI communicates directly with the Wazuh Indexer (OpenSearch), not the Manager API. OpenSearch uses its own internal security plugin (Basic Auth), which requires appropriate permissions.

### Credential Fallback
If you configure `indexer_url` but omit `indexer_user` and `indexer_password`, the CLI automatically falls back to using the Manager API `user` and `password`. This works well for default installations where both services use `admin` credentials.

### Indexer Permissions Setup

If you use a specific API user (like `wazuh-cli`) and encounter a `403 Forbidden` error during alert queries, your CLI user lacks the necessary read permissions on the OpenSearch indices (`wazuh-alerts-*`).

To fix this, you must map your user to a role with proper index permissions:

1. Log into the **Wazuh/OpenSearch Dashboards** using an admin account.
2. Go to **Security** -> **Roles** and click **Create role** (e.g., `wazuh_cli_reader`).
3. Under **Index permissions**:
   - **Index**: `wazuh-alerts-*`
   - **Permissions**: Add `indices:data/read/search` (or select the `read` action group).
4. Save the role.
5. Go to **Mapped users** for this role and add your CLI user (e.g., `wazuh-cli`).

Alternatively, you can bypass role mapping by explicitly configuring the CLI to use an admin indexer account while keeping the Manager API account restricted:

```bash
wazuh-cli config set indexer_user admin
wazuh-cli config set indexer_password "admin_password"
```
