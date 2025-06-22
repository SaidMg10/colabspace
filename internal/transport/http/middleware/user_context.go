package middleware

import (
	"errors"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/gin-gonic/gin"
)

// Se crea una funcion que obtenga el id desde el contexto
// mejor dicho que obtenga el user desde el contexto
func GetUserFromContext(c *gin.Context) (*user.User, error) {
	// Extraer el usuario autenticado guardado en el contexto de gin
	usr, exists := c.Get(userKey)
	if !exists {
		return nil, errors.New("user not exists in the context")
	}
	// Se realiza un type assertion para asegurarse de que el valor tiene el tipo necesario de dato
	u, ok := usr.(*user.User)
	if !ok {
		return nil, errors.New("invalid user type in contex")
	}
	// Retorna el u validado
	return u, nil
}
