package keyring

import (
	"github.com/zalando/go-keyring"
	"go.bbkane.com/namedenv/domain"
)

type OSKeyring struct {
	service string
}

func (k *OSKeyring) Set(key string, value string) error {
	return keyring.Set(k.service, key, value)
}

func (k *OSKeyring) Get(key string) (string, error) {
	return keyring.Get(k.service, key)
}

func (k *OSKeyring) Delete(key string) error {
	return keyring.Delete(k.service, key)
}

func (k *OSKeyring) Service() string {
	return k.service
}

func NewOSKeyring(service string) domain.Keyring {
	return &OSKeyring{
		service: service,
	}
}
