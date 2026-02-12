package response

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AllUsers struct {
	Entries []UserResponse `json:"entries"`
}

func BuildUserResponse(id uint, name, email, role string) UserResponse {
	return UserResponse{
		ID: id,
		Name: name,
		Email: email,
		Role: role,
	}
}

func BuildAllUser(users []UserResponse) AllUsers{
	return AllUsers{
		Entries: users,
	}
}