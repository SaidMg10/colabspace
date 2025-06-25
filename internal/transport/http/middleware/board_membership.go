package middleware

import (
	"net/http"
	"strconv"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CheckBoardMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		usr, exists := c.Get(userKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		user, ok := usr.(*user.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type"})
			c.Abort()
			return
		}
		userID := user.ID

		boardID, err := strconv.ParseInt(c.Param("boardID"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		// Verificar si es dueño
		board, err := m.Services.Board.GetById(ctx, boardID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
			c.Abort()
			return
		}
		if board.UserID == userID {
			c.Next()
			return
		}

		// Verificar si es miembro
		isMember, err := m.Services.BoardUsers.IsMember(ctx, boardID, userID)
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

		c.Next()
	}
}
