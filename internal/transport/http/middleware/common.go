package middleware

import (
	"github.com/SaidMg10/colabspace/internal/app"
	"github.com/SaidMg10/colabspace/internal/app/auth"
)

type Middleware struct {
	Authenticator auth.Authenticator
	Services      *app.Services
}

func NewMiddleware(authenticator auth.Authenticator, services *app.Services) *Middleware {
	return &Middleware{
		Authenticator: authenticator,
		Services:      services,
	}
}
