path "auth/token/*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}

path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

path "auth/approle/role/go_vault_session_role/role-id" {
  capabilities = [ "read" ]
}

path "auth/approle/role/go_vault_session_role/secret-id" {
  capabilities = ["create", "read", "update"]
}

path "secret/*" {
  capabilities = ["create", "update", "read"]
}