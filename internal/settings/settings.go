package settings

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Settings struct {
	storage Storage
}

func New(storage Storage) (*Settings, error) {
	return &Settings{
		storage: storage,
	}, nil
}

func (s *Settings) GetOrGenerate(name string, size uint) (string, error) {
	value, err := s.storage.Get(name)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if value == "" {
		// Generate a new secret
		value, err = generateRandomString(size)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		err = s.storage.Set(name, value)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}
	return value, nil
}

func generateRandomString(n uint) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"
	ret := make([]byte, n)
	for i := uint(0); i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
