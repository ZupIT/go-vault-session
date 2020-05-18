package main

import (
	"github.com/ernelio/go-vault-session/pkg/token"
	"github.com/hashicorp/vault/api"
	"log"

	"github.com/ernelio/go-vault-session/pkg/login"
)

func main() {
	client := vaultConfig()
	vaultStarter(client)
}

func vaultConfig() *api.Client {
	vaultConfig := api.DefaultConfig()
	if err := vaultConfig.ReadEnvironment(); err != nil {
		log.Fatal(err)
	}

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func vaultStarter(client *api.Client) {
	vaultLogin := login.NewHandler(client)
	secret := vaultLogin.Handle()
	ch := make(chan string)

	renewal := token.NewHandler(client, secret, ch)
	renewal.Handle()
}
