package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SaidMg10/colabspace/internal/app/board"
	"github.com/SaidMg10/colabspace/internal/transport/http/middleware"
	"github.com/SaidMg10/colabspace/internal/utils"
	"github.com/SaidMg10/colabspace/internal/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BoardHandler struct {
	Service    *board.Service
	Logger     *zap.SugaredLogger
	Middleware *middleware.Middleware
}

func NewBoardHandler(service *board.Service, logger *zap.SugaredLogger, mw *middleware.Middleware) *BoardHandler {
	return &BoardHandler{
		Service:    service,
		Logger:     logger,
		Middleware: mw,
	}
}

// RegisterRoutes registra las rutas HTTP relacionadas con usuarios dentro del grupo de rutas pasado como parámetro.
// Agrupa rutas bajo el prefijo /boards y asigna los handlers correspondientes a cada método HTTP.
func (h *BoardHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Creamos un subgrupo /boards para organizar las rutas relacionadas con usuarios.
	boards := r.Group("/boards")
	{
		// Ruta POST /boards/ para crear un nuevo usuario.
		// El método Create es el handler que procesará esta ruta.
		boards.POST("/", h.Create)
		boards.GET("/", h.Get)
		boards.GET("/:id", h.GetById)
		boards.PATCH("/:id", h.Middleware.CheckBoardOwnership(), h.Update)
		boards.DELETE("/:id", h.Middleware.CheckBoardOwnership(), h.Delete)
	}
}

func (h *BoardHandler) Create(c *gin.Context) {
	var cBR board.CreateBoardRequest
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error has ocurred"})
		return
	}
	cBR.UserID = user.ID
	if err := c.ShouldBindJSON(&cBR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := validation.Validate.Struct(cBR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	b, err := h.Service.Create(ctx, cBR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error has ocurred"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"board": b,
	})
}

func (h *BoardHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	b, err := h.Service.Get(ctx)
	if err != nil {
		h.Logger.Errorf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if len(b) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": b,
	})
}

func (h *BoardHandler) GetById(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ctx := c.Request.Context()
	b, err := h.Service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
			return
		}
		h.Logger.Errorf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if b == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func (h *BoardHandler) Update(c *gin.Context) {
	// Obtenemos el parametro del id
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	// Validacion del payload pasado
	var uBR board.UpdateBoardRequest
	if err := c.ShouldBindJSON(&uBR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := validation.Validate.Struct(uBR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Sacamos el context
	ctx := c.Request.Context()
	// Pasamos el payload y el id a la funcion del service
	b, err := h.Service.Update(ctx, id, uBR)
	if err != nil {
		h.Logger.Errorf("Error deleting board: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"board": b,
	})
}

func (h *BoardHandler) Delete(c *gin.Context) {
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
		h.Logger.Errorf("Error deleting a board: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Board was deleted"})
}
