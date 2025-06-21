package handlers

import (
	"errors"
	"net/http"

	"github.com/SaidMg10/colabspace/internal/app/auth"
	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/SaidMg10/colabspace/internal/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Service *auth.Service
	Logger  *zap.SugaredLogger
}

func NewAuthHandler(service *auth.Service, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Logger:  logger,
	}
}

// RegisterRoutes registra las rutas HTTP relacionadas con usuarios dentro del grupo de rutas pasado como parámetro.
// Agrupa rutas bajo el prefijo /users y asigna los handlers correspondientes a cada método HTTP.
func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Creamos un subgrupo /users para organizar las rutas relacionadas con usuarios.
	users := r.Group("/auth")
	{
		// Ruta POST /users/ para crear un nuevo usuario.
		// El método Create es el handler que procesará esta ruta.
		users.POST("/", h.Login)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload auth.Login
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&payload); err != nil {
		h.Logger.Warnf("Login payload invalid: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	token, err := h.Service.Login(ctx, payload)
	if err != nil {
		h.Logger.Warnf("Login failed for email %s: %v", payload.Email, err)
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *UserHandler) Register(c *gin.Context) {
	var cUR user.CreateUserRequest
	if err := c.ShouldBindJSON(&cUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := validation.Validate.Struct(cUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()

	u, err := h.Service.Create(ctx, cUR)
	if err != nil {
		h.Logger.Errorf("Error creating user: %v", err)
		switch {
		case errors.Is(err, user.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		case errors.Is(err, user.ErrDuplicateUsername):
			c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"user": u,
	})
}
