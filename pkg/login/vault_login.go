package login

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

var (
	authType                                 = os.Getenv("VAULT_AUTHENTICATION")
	roleId                                   = os.Getenv("VAULT_ROLE_ID")
	secretId                                 = os.Getenv("VAULT_SECRET_ID")
	k8sRole                                  = os.Getenv("VAULT_K8S_ROLE")
	k8sPath                                  = getEnv("VAULT_K8S_PATH", "auth/kubernetes/login")
	k8sAuth                                  = getEnv("VAULT_K8S_AUTH", "KUBERNETES")
	appRolePath                              = getEnv("VAULT_APP_ROLE_PATH", "auth/approle/login")
	appRoleAuth                              = getEnv("VAULT_APP_ROLE_AUTH", "APPROLE")
	defaultKubernetesServiceAccountTokenFile = getEnv("VAULT_DEFAULT_K8S_SERVICE_ACCOUNT_TOKEN_FILE", "/var/run/secrets/kubernetes.io/serviceaccount/token")
)

type Handler interface {
	Handle() *api.Secret
}

type Manager struct {
	client *api.Client
}

func NewHandler(c *api.Client) *Manager {
	return &Manager{client: c}
}

func (l *Manager) Handle() *api.Secret {
	path, body := authResolver()
	res, err := l.client.Logical().Write(path, body)
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Setenv(api.EnvVaultToken, res.Auth.ClientToken)

	return res
}

func authResolver() (string, map[string]interface{}) {
	var authPath string
	var data map[string]interface{}

	switch authType {
	case appRoleAuth:
		authPath = appRolePath
		data = map[string]interface{}{
			"role_id":   roleId,
			"secret_id": secretId,
		}
	case k8sAuth:
		authPath = k8sPath
		data = map[string]interface{}{
			"jwt":  readK8sJwt(),
			"role": k8sRole,
		}
	}

	return authPath, data
}

func readK8sJwt() string {
	jwt, err := ioutil.ReadFile(defaultKubernetesServiceAccountTokenFile)
	if err != nil {
		log.Fatal(err)
	}

	return string(jwt)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
