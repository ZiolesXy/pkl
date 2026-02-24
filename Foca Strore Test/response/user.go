package response

type UserMiniResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AllUsers struct {
	Entries []UserResponse `json:"entries"`
}

type UserProfileResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Role            string `json:"role"`
	Address         []AddressResponse `json:"address"`
}

func BuildUserResponse(id uint, name, email, role string) UserResponse {
	return UserResponse{
		ID:    id,
		Name:  name,
		Email: email,
		Role:  role,
	}
}

func BuildAllUser(users []UserResponse) AllUsers {
	return AllUsers{
		Entries: users,
	}
}

func BuildUserProfileResponse(id uint, name, email, profileImageURL, role string, address []AddressResponse) UserProfileResponse {
	return UserProfileResponse{
		ID:              id,
		Name:            name,
		Email:           email,
		ProfileImageURL: profileImageURL,
		Role:            role,
		Address:         address,
	}
}
