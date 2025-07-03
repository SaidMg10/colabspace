package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/SaidMg10/colabspace/internal/utils"
	"github.com/SaidMg10/colabspace/internal/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler es el struct que contiene el service asociado a la lógica de usuario.
// Permite que los handlers tengan acceso a las funciones del servicio para procesar la lógica de negocio.
type UserHandler struct {
	Service *user.Service
	Logger  *zap.SugaredLogger
}

// NewUserHandler crea una nueva instancia de UserHandler recibiendo el service correspondiente.
// Esto permite inyectar dependencias y facilita pruebas unitarias.
func NewUserHandler(service *user.Service, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		Service: service,
		Logger:  logger,
	}
}

// RegisterRoutes registra las rutas HTTP relacionadas con usuarios dentro del grupo de rutas pasado como parámetro.
// Agrupa rutas bajo el prefijo /users y asigna los handlers correspondientes a cada método HTTP.
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Creamos un subgrupo /users para organizar las rutas relacionadas con usuarios.
	users := r.Group("/users")
	{
		// Ruta POST /users/ para crear un nuevo usuario.
		// El método Create es el handler que procesará esta ruta.
		users.POST("/", h.Create)
		users.GET("/", h.Get)
		users.GET("/:id", h.GetById)
		users.PATCH("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
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

	u, err := h.Service.Create(ctx, &cUR)
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

func (h *UserHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	u, err := h.Service.Get(ctx)
	if err != nil {
		h.Logger.Errorf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if len(u) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

func (h *UserHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	ctx := c.Request.Context()
	u, err := h.Service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.Logger.Errorf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

func (h *UserHandler) Update(c *gin.Context) {
	// Obtenemos el parametro del id
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	// Validacion del payload pasado
	var uUR user.UpdateUserRequest
	if err := c.ShouldBindJSON(&uUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := validation.Validate.Struct(uUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Sacamos el context
	ctx := c.Request.Context()
	// Pasamos el payload y el id a la funcion del service
	u, err := h.Service.Update(ctx, id, &uUR)
	if err != nil {
		h.Logger.Errorf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"user": u,
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	// Obtenemos el parametro
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	// Sacamos el context
	ctx := c.Request.Context()
	// Pasamos lo requerido para la funcion del service
	if err := h.Service.Delete(ctx, id); err != nil {
		h.Logger.Errorf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User was deleted"})
}
