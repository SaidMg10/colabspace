package models

import "golang.org/x/crypto/bcrypt"

type Password struct {
	hash []byte
}

func (p *Password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	return nil
}

func (p *Password) Compare(plain string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
}

// Getter para obtener el hash (para guardar en la DB)
func (p *Password) Hash() []byte {
	return p.hash
}
