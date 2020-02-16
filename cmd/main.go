package main

import (
	"github.com/hashicorp/vault/api"

	"github.com/ZupIT/go-vault-session/pkg/login"
)

func main() {
	client := vaultConfig()
	vaultStarter(client)
}

func vaultConfig() *api.Client {
	vaultConfig := api.DefaultConfig()
	_ = vaultConfig.ReadEnvironment()
	client, _ := api.NewClient(vaultConfig)
	return client
}

func vaultStarter(client *api.Client) {
	vaultAuth := login.NewHandler(client)
	_ = vaultAuth.HandleLogin()
}
