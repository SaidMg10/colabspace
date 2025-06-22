package middleware

import (
	"net/http"
	"strconv"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CheckBoardOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer el user ID del contexto
		usr, exists := c.Get(userKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		user, ok := usr.(*user.User) // o el tipo que estés usando
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type"})
			c.Abort()
			return
		}
		// Extraer el user id
		userID := user.ID
		// Extraer el board id del path param
		boardID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		board, err := m.Services.Board.GetById(ctx, boardID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "board not found with this id"})
			c.Abort()
			return
		}
		// Verificar propiedad
		if board.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: you are not the owner"})
			c.Abort()
			return
		}

		// Continuar
		c.Next()
	}
}
