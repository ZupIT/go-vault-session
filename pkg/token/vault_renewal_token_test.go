package token

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/go-vault-session/pkg/login"
)

func Test_Should_renewal_vault_token_success(t *testing.T) {
	client := vaultClient()
	l := login.NewHandler(client)
	secret := l.HandleLogin()

	renewal := NewRenewalHandler(client, secret)
	renewal.HandleRenewal()

	vaultToken := os.Getenv(api.EnvVaultToken)
	assert.NotEmpty(t, vaultToken)
	client.SetToken(vaultToken)

	pathWithKey := "secret/data/my-secret"
	body := map[string]interface{}{
		"data": map[string]string{"test": "test_ok"},
	}
	_, _ = client.Logical().Write(pathWithKey, body)

	res, err := client.Logical().Read(pathWithKey)
	assert.Nil(t, err)

	data := res.Data["data"].(map[string]interface{})
	nameTest := data["test"]
	assert.Equal(t, "test_ok", nameTest)
}

func vaultClient() *api.Client {
	vaultConfig := api.DefaultConfig()
	_ = vaultConfig.ReadEnvironment()
	client, _ := api.NewClient(vaultConfig)
	return client
}
