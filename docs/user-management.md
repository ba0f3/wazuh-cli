# API User Management

This guide explains how to create and configure new API users in Wazuh.

## Creating a New User

There are two primary ways to create a new user: via the Wazuh Dashboard or the API (cURL).

### Method 1: Wazuh Dashboard (Easiest)

1. Log in to the **Wazuh Dashboard** with an administrator account (default is `admin`).
2. Open the main menu (☰) in the top-left corner.
3. Navigate to **Server management** > **Security**.
4. Select the **Users** tab and click the **Create user** button.
5. Enter the user details:
   - **Username**: e.g., `api_service_user`
   - **Password**: Must comply with the password policy (8-64 characters, uppercase, lowercase, numbers, and symbols).
6. Click **Create** to finish.

### Method 2: Wazuh API (cURL)

First, obtain a JWT token using an existing administrator account:

```bash
# 1. Get Token (replace 'admin:secret' with your credentials)
TOKEN=$(curl -u admin:secret -k -X POST "https://<WAZUH_MANAGER_IP>:55000/security/user/authenticate?raw=true")

# 2. Create the new user
curl -k -X POST "https://<WAZUH_MANAGER_IP>:55000/security/users" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"username": "new_api_user", "password": "YourPassword123!"}'
```

---

## Important Next Steps: Assigning Roles

Creating a user is not enough; you must assign them **Roles** so they have permissions to perform actions (like reading logs or managing agents).

1. Go to **Server management** > **Security**.
2. Select the **Roles mapping** tab.
3. Link your new user to an existing role:
   - `administrator`: Full access.
   - `readonly`: Read-only access.
   - Or a custom role with specific permissions.

### Enabling "Run As"
If you intend to use this account to perform actions on behalf of other users, ensure that the `run_as: true` option is enabled in the Wazuh Dashboard configuration file:
`/usr/share/wazuh-dashboard/data/wazuh/config/wazuh.yml`

## Using the User with wazuh-cli

Once the user is created and assigned a role, you can configure `wazuh-cli` to use it:

```bash
./bin/wazuh-cli auth login --url https://<WAZUH_MANAGER_IP>:55000 --user new_api_user
```
