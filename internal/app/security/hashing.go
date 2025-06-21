package security

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Hash = hash

	return nil
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
}
