package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	Hasher struct {
		cost int
	}

	Interactor interface {
		HashPassword(string) (string, error)
		CompareAndHash(hashed string, curr string) bool
	}
)

func (h *Hasher) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, h.cost)

	return string(hashedPasswordBytes), err
}

func (h *Hasher) CompareAndHash(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

var _ Interactor = (*Hasher)(nil)

func New(cost int) *Hasher {
	return &Hasher{
		cost: cost,
	}
}
