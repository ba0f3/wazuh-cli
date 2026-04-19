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
