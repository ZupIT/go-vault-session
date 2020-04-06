package token

import (
	"log"

	"github.com/hashicorp/vault/api"
)

type Handler interface {
	Handle()
}

type Manager struct {
	client *api.Client
	secret *api.Secret
}

func NewHandler(c *api.Client, s *api.Secret) *Manager {
	return &Manager{client: c, secret: s}
}

func (c *Manager) Handle() {
	r, _ := c.client.NewRenewer(&api.RenewerInput{Secret: c.secret})

	go func() {
		go r.Renew()
		defer r.Stop()
		for {
			select {
			case err := <-r.DoneCh():
				if err != nil {
					log.Fatal(err)
				}
			case _ = <-r.RenewCh():
				log.Println("Token successfully renewed")
			}
		}
	}()
}
