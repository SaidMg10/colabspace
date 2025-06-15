package handlers

import "github.com/SaidMg10/colabspace/internal/server"

type Handler struct {
	App *server.Application
}

func NewHandler(app *server.Application) *Handler {
	return &Handler{App: app}
}
