package login

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func Test_Should_login_in_vault_with_success(t *testing.T) {
	assert.NotEmpty(t, roleId)
	assert.NotEmpty(t, secretId)
	assert.Equal(t, appRoleAuth, authType)

	client := config()
	login := NewHandler(client)
	_ = login.Handle()

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

func config() *api.Client {
	vaultConfig := api.DefaultConfig()
	_ = vaultConfig.ReadEnvironment()
	client, _ := api.NewClient(vaultConfig)
	return client
}
