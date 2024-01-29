package keyring

import (
	"errors"
	"runtime"

	"github.com/zalando/go-keyring"
	"go.bbkane.com/namedenv/domain"
)

type OSKeyring struct {
	service string
}

func wrapErr(err error) error {
	if err == nil {
		return nil
	} else if errors.Is(err, keyring.ErrNotFound) {
		return domain.ErrNotFound
	} else if errors.Is(err, keyring.ErrNotFound) {
		return domain.ErrNotFound
	} else {
		return err
	}
}

func (k *OSKeyring) Set(key string, value string) error {
	err := keyring.Set(k.service, key, value)
	return wrapErr(err)
}

func (k *OSKeyring) Get(key string) (string, error) {
	res, err := keyring.Get(k.service, key)
	err = wrapErr(err)
	return res, err
}

func (k *OSKeyring) Delete(key string) error {
	err := keyring.Delete(k.service, key)
	return wrapErr(err)
}

func (k *OSKeyring) Service() string {
	return k.service
}

func NewOSKeyring(service string) domain.Keyring {

	os := runtime.GOOS

	if os != "darwin" && os != "linux" && os != "windows" {
		panic("unsupported OS: " + os)
	}

	return &OSKeyring{
		service: service,
	}
}
