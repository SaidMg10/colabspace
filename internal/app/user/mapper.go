package user

import "github.com/SaidMg10/colabspace/internal/app/security"

func MapCreateUserRequestToUser(req CreateUserRequest) (*User, error) {
	var pwd security.Password
	if err := pwd.Set(req.Password); err != nil {
		return nil, err
	}

	user := &User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  pwd,
		Role:      RoleUser,
	}

	return user, nil
}

func MapUpdateUserRequestToUser(req UpdateUserRequest, userExist *User) (*User, error) {
	// Primero copiamos los datos del user existente
	updated := *userExist
	// Empezamos con la validacion manual de cada campo
	// En Golang (y al menos en mi noción)
	// Los campos al ser trabajados como patch se deben ir editando uno por uno
	// Esto hace que cada validación deba ser realizada manualmente
	// Esto para tener una persistencia correcta entre los campos en la bd
	if req.FirstName != nil {
		updated.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		updated.LastName = *req.LastName
	}
	// En este campo hacemos algo especial, ya que el tipo Password
	// no acepta directamente un string plano. Requiere usar el método
	// Set(string) para generar internamente el hash antes de asignarlo.
	// Además, verificamos que el campo no solo no sea nil, sino que también
	// no esté vacío (""), para evitar sobreescribir con valores inválidos.
	if req.Password != nil && *req.Password != "" {
		// Generamos el hash de la nueva contraseña usando el método Set,
		// similar al proceso en la creación del usuario. Aquí usamos el puntero
		// req.Password para acceder directamente al valor enviado en la request.
		var pwd security.Password
		if err := pwd.Set(*req.Password); err != nil {
			return nil, err
		}
		updated.Password = pwd
	}
	// Devolvemos el User preparado para ser enviado a la bd
	return &updated, nil
}

func MapToUsersSummaryResponse(u []User) []UserSummaryResponse {
	resp := make([]UserSummaryResponse, len(u))
	for i, v := range u {
		resp[i] = UserSummaryResponse{
			ID:        v.ID,
			Username:  v.Username,
			FirstName: v.FirstName,
			LastName:  v.LastName,
		}
	}
	return resp
}

func MapToUserSummaryResponse(u *User) *UserSummaryResponse {
	resp := &UserSummaryResponse{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
	return resp
}
