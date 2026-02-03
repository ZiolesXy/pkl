package request

type UserPost struct {
	Name string `json:"name" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

type UserPut struct {
	Name *string `json:"name"`
	RoleID *uint `json:"role_id,omitempty"`
}