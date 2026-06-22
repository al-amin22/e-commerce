package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID          any    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        any    `json:"role"`
	IsVerified  bool   `json:"is_verified,omitempty"`
}
