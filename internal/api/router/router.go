package router

import (
	"github.com/SaidMg10/colabspace/internal/api/handlers"
	"github.com/SaidMg10/colabspace/internal/server"
	"github.com/gin-gonic/gin"
)

func Mount(app *server.Application) *gin.Engine {
	r := gin.New()

	h := handlers.NewHandler(app)

	r.GET("/health", h.HealthCheck)

	users := r.Group("/user")
	{
		users.POST("", h.CreateUser)
	}

	return r
}
