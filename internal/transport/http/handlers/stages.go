package handlers

import (
	"net/http"
	"strconv"

	"github.com/SaidMg10/colabspace/internal/app/stage"
	"github.com/SaidMg10/colabspace/internal/transport/http/middleware"
	"github.com/SaidMg10/colabspace/internal/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler es el struct que contiene el service asociado a la lógica de usuario.
// Permite que los handlers tengan acceso a las funciones del servicio para procesar la lógica de negocio.
type StageHandler struct {
	Service    *stage.Service
	Logger     *zap.SugaredLogger
	Middleware *middleware.Middleware
}

// NewUserHandler crea una nueva instancia de UserHandler recibiendo el service correspondiente.
// Esto permite inyectar dependencias y facilita pruebas unitarias.
func NewStageHandler(
	service *stage.Service,
	logger *zap.SugaredLogger,
	middleware *middleware.Middleware,
) *StageHandler {
	return &StageHandler{
		Service:    service,
		Logger:     logger,
		Middleware: middleware,
	}
}

// RegisterRoutes registra las rutas HTTP relacionadas con usuarios dentro del grupo de rutas pasado como parámetro.
// Agrupa rutas bajo el prefijo /users y asigna los handlers correspondientes a cada método HTTP.
func (h *StageHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Creamos un subgrupo /users para organizar las rutas relacionadas con usuarios.
	boards := r.Group("/boards")
	{
		stages := boards.Group("/:boardID/stages")
		{
			stages.POST("/", h.Create)
			stages.GET("/public", h.GetPublic)
			stages.GET("/:stageID", h.GetByIDPublic)
			stages.PATCH("/:stageID", h.Update)
			stages.DELETE("/:stageID", h.Delete)
		}
	}
}

func (h *StageHandler) Create(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	var cSR stage.CreateStageRequest
	if err := c.ShouldBindJSON(&cSR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := validation.Validate.Struct(cSR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cSR.BoardId = boardID
	ctx := c.Request.Context()

	st, err := h.Service.Create(ctx, &cSR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating user",
		})
		h.Logger.Errorf("Error creating user: %v", err)
	}
	c.JSON(http.StatusCreated, gin.H{
		"stage": st,
	})
}

func (h *StageHandler) GetPublic(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid boardID",
		})
	}
	ctx := c.Request.Context()
	st, err := h.Service.Get(ctx, boardID)
	if err != nil {
		h.Logger.Errorf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if len(st) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"stages": st,
	})
}

func (h *StageHandler) GetByIDPublic(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	stageIDStr := c.Param("stageID")
	stageID, err := strconv.ParseInt(stageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	ctx := c.Request.Context()
	st, err := h.Service.GetById(ctx, boardID, stageID)
	if err != nil {
		h.Logger.Errorf("Error getting stage: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"stage": st,
	})
}

func (h *StageHandler) Update(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	stageIDStr := c.Param("stageID")
	stageID, err := strconv.ParseInt(stageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stage id"})
		return
	}

	var uSR stage.UpdateStageRequest
	if err := c.ShouldBindJSON(&uSR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validar estructura
	if err := validation.Validate.Struct(uSR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	resp, err := h.Service.Update(ctx, boardID, stageID, &uSR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stage": resp,
	})
}

func (h *StageHandler) Delete(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid boardID",
		})
		return
	}
	stageIDStr := c.Param("stageID")
	stageID, err := strconv.ParseInt(stageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid stageID",
		})
		return
	}
	ctx := c.Request.Context()
	if err := h.Service.Delete(ctx, boardID, stageID); err != nil {
		h.Logger.Errorf("Error ocurred while delete a stage %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
