package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type bcryptProvider struct {
	cost int
}

func NewBcryptProvider(cost int) Provider {
	return &bcryptProvider{
		cost: cost,
	}
}

func (p *bcryptProvider) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (p *bcryptProvider) CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
