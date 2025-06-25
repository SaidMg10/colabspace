package handlers

import (
	"net/http"
	"strconv"

	boardusers "github.com/SaidMg10/colabspace/internal/app/board_users"
	"github.com/SaidMg10/colabspace/internal/transport/http/middleware"
	"github.com/SaidMg10/colabspace/internal/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BoardUsersHandler struct {
	Service    *boardusers.Service
	Logger     *zap.SugaredLogger
	Middleware *middleware.Middleware
}

func NewBoardUsersHandler(service *boardusers.Service, logger *zap.SugaredLogger, mw *middleware.Middleware) *BoardUsersHandler {
	return &BoardUsersHandler{
		Service:    service,
		Logger:     logger,
		Middleware: mw,
	}
}

func (h *BoardUsersHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Creamos un subgrupo /users para organizar las rutas relacionadas con usuarios.
	boardUsers := r.Group("/boards")
	{
		// Ruta POST /users/ para crear un nuevo usuario.
		// El método Create es el handler que procesará esta ruta.
		boardUsers.POST("/:boardID/members", h.Middleware.CheckBoardOwnership(), h.Add)
		boardUsers.GET("/:boardID/members", h.Middleware.CheckBoardMembership(), h.Get)
		boardUsers.DELETE("/:boardID/members/leave", h.Middleware.CheckUserMembership(), h.Leave)
		boardUsers.DELETE("/:boardID/members/:userID", h.Middleware.CheckBoardOwnership(), h.Delete)
	}
}

func (h *BoardUsersHandler) Add(c *gin.Context) {
	var cBUR boardusers.CreateBoardUsersRequest

	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	if err := c.ShouldBindJSON(&cBUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	cBUR.BoardID = boardID

	if err := validation.Validate.Struct(cBUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.Service.Add(ctx, &cBUR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *BoardUsersHandler) Get(c *gin.Context) {
	boardIdStr := c.Param("boardID")
	id, err := strconv.ParseInt(boardIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	ctx := c.Request.Context()
	bUsrs, err := h.Service.GetMembers(ctx, id)
	if err != nil {
		h.Logger.Errorf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if len(bUsrs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"members": bUsrs,
	})
}

func (h *BoardUsersHandler) Leave(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	id, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error has ocurred"})
		return
	}
	ctx := c.Request.Context()
	if err := h.Service.Leave(ctx, id, user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error leaving the board"})
	}
}

func (h *BoardUsersHandler) Delete(c *gin.Context) {
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}
	ctx := c.Request.Context()
	if err := h.Service.Leave(ctx, boardID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error leaving the board"})
	}
}
