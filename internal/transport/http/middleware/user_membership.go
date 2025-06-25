package middleware

import (
	"net/http"
	"strconv"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CheckUserMembership() gin.HandlerFunc {
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
		boardID, err := strconv.ParseInt(c.Param("boardID"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		// Verificar si es owner
		board, err := m.Services.Board.GetById(ctx, boardID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "board not found with this id"})
			c.Abort()
			return
		}

		if board.UserID == userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: you're the owner"})
			c.Abort()
			return
		}

		isMember, err := m.Services.BoardUsers.Repo.IsMember(ctx, boardID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking membership"})
			c.Abort()
			return
		}
		if !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not a member"})
			c.Abort()
			return
		}

		// Continuar
		c.Next()
	}
}
