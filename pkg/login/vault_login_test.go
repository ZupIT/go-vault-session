package login

import (
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
)

var client *api.Client

func TestMain(m *testing.M) {
	client = config()
	os.Exit(m.Run())
}

func TestHandle(t *testing.T) {

	type in struct {
		path string
		data map[string]string
	}

	type out struct {
		err  error
		want map[string]interface{}
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "login success",
			in: &in{
				path: "secret/data/my-secret",
				data: map[string]string{"test": "test_ok"},
			},
			out: &out{
				err:  nil,
				want: map[string]interface{}{"test": "test_ok"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login := NewHandler(client)
			login.Handle()

			vaultToken := os.Getenv(api.EnvVaultToken)

			if vaultToken == "" {
				t.Errorf("Login(%s) must be a vault token valid", tt.name)
			}

			client.SetToken(vaultToken)

			body := map[string]interface{}{
				"data": tt.in.data,
			}
			_, _ = client.Logical().Write(tt.in.path, body)

			res, err := client.Logical().Read(tt.in.path)

			if err != tt.out.err {
				t.Errorf("Login(%s) got %v, want %v", tt.name, err, tt.out.err)
			}

			got := res.Data["data"].(map[string]interface{})

			if !reflect.DeepEqual(tt.out.want, got) {
				t.Errorf("Login(%s) got %v, want %v", tt.name, got, tt.out.want)
			}
		})
	}
}

func config() *api.Client {
	vaultConfig := api.DefaultConfig()
	_ = vaultConfig.ReadEnvironment()
	client, _ := api.NewClient(vaultConfig)
	return client
}
