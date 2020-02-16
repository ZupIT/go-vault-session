package login

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

var (
	authType = os.Getenv("VAULT_AUTHENTICATION")
	roleId   = os.Getenv("VAULT_ROLE_ID")
	secretId = os.Getenv("VAULT_SECRET_ID")
	k8sRole  = os.Getenv("VAULT_K8S_ROLE")
)

const (
	k8sPath                                  = "auth/kubernetes/login"
	k8sAuth                                  = "KUBERNETES"
	appRolePath                              = "auth/approle/login"
	appRoleAuth                              = "APPROLE"
	defaultKubernetesServiceAccountTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type Login struct {
	client *api.Client
}

func NewHandler(c *api.Client) *Login {
	return &Login{client: c}
}

func (l *Login) HandleLogin() *api.Secret {
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
		log.Panicln(err)
	}

	return string(jwt)
}
