package users

type UserResponse struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func NewUserResponse(user *User) UserResponse {
	return UserResponse{
		UUID:  user.UUID.String(),
		Name:  user.Name,
		Email: user.Email,
	}
}
