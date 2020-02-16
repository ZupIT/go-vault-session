package token

import (
	"log"

	"github.com/hashicorp/vault/api"
)

type Renewal struct {
	client *api.Client
	secret *api.Secret
}

func NewRenewalHandler(c *api.Client, s *api.Secret) *Renewal {
	return &Renewal{client: c, secret: s}
}

func (c *Renewal) HandleRenewal() {
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
