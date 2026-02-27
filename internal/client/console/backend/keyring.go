package backend

import (
	"github.com/zalando/go-keyring"
)

const keyringService = "connectfour"
const keyringAccount = "jwt"

func saveToKeyring(token string) error {
	return keyring.Set(keyringService, keyringAccount, token)
}

func loadFromKeyring() (string, error) {
	return keyring.Get(keyringService, keyringAccount)
}

func removeFromKeyring() error {
	return keyring.Delete(keyringService, keyringAccount)
}
