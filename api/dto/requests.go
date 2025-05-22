package dto

// CreateUserRequest represents the create user request
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    uint   `json:"role_id"`
}

// UpdateUserRequest represents the update user request
type UpdateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    uint   `json:"role_id"`
}

// RegisterRequest represents the register request
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    uint   `json:"role_id"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse represents the register response
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
}

// CreateBlogRequest represents the create blog request
type CreateBlogRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Published bool   `json:"published"`
}

// UpdateBlogRequest represents the update blog request
type UpdateBlogRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Published bool   `json:"published"`
}
