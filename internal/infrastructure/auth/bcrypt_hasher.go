package auth

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	return &BcryptPasswordHasher{
		cost: cost,
	}
}

func (b *BcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (b *BcryptPasswordHasher) ComparePassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
