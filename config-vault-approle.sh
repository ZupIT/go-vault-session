#!/bin/sh

C_DIR=$1
VAULT_HOST=$2

mkdir -p /tmp/vault

echo "Exporting vault vars..."
export VAULT_TOKEN="c8159b3f-e7fd-4be4-badd-cd0b78207381"
export VAULT_ADDR=$VAULT_HOST

echo "Installing vault cli 1.3.0..."
rm -rf /tmp/vault/vault
unzip "$C_DIR"/resources/vault_1.3.0_"$(uname -s)"_amd64.zip -d /tmp/vault/

# Create a vault policy named vault_session_policy
/tmp/vault/vault policy write vault_session_policy "$C_DIR"/resources/vault_session_policy.hcl

# Enable auth method APPROLE
/tmp/vault/vault auth enable approle

# Add vault_session_policy for APPROLE authentication with 15-second token period
/tmp/vault/vault write auth/approle/role/go_vault_session_role policies=vault_session_policy period=15s

rm -f /tmp/vault/role-id.txt
rm -f /tmp/vault/secret-id.txt

role_response=$(/tmp/vault/vault read -format=json auth/approle/role/go_vault_session_role/role-id)
echo "role_response $role_response"
role_id=$(echo "$role_response" | "$C_DIR/jq-$(uname -s)" -j '.data.role_id')
echo "role_id: $role_id"
eval echo "$role_id" >>/tmp/vault/role-id.txt

secret_response=$(/tmp/vault/vault write -force -format=json auth/approle/role/go_vault_session_role/secret-id)
echo "secret_response: $secret_response"
secret_id=$(echo "$secret_response" | "$C_DIR/jq-$(uname -s)" -j '.data.secret_id')
echo "secret_id: $secret_id"
eval echo "$secret_id" >>/tmp/vault/secret-id.txt

unset VAULT_TOKEN
