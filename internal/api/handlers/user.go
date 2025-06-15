package handlers

import (
	"log"
	"net/http"

	"github.com/SaidMg10/colabspace/internal/dto"
	"github.com/SaidMg10/colabspace/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var userDto dto.CreateUserDto
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()

	var pwd models.Password
	err := pwd.Set(userDto.Password)
	if err != nil {
		// Por ejemplo, responder con un error 500 (Internal Server Error)
		log.Printf("Error al generar hash de la contraseña: %v", err)
		// Si estás en un handler HTTP con Gin, por ejemplo:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar la contraseña"})
		return
	}

	user := &models.User{
		Username:  userDto.Username,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Email:     userDto.Email,
		Password:  pwd,
		Role:      models.RoleUser,
	}

	if err := h.App.Store.Users.Create(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	resp := dto.CreateUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"user": resp})
}
