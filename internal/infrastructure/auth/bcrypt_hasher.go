package auth

import "golang.org/x/crypto/bcrypt"

type BcryptPassworHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPassworHasher {
	return &BcryptPassworHasher{
		cost: cost,
	}
}

func (b *BcryptPassworHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (b *BcryptPassworHasher) ComparePassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
