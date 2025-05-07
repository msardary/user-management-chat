package auth

import "time"


type User struct {
    ID           int64
    Username     string
    PasswordHash string
	FirstName    string
	LastName     string
	Mobile       string
	Email        string
    IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
    ID         int
    UserID     int
    TokenHash  string
    ExpiresAt  time.Time
    Revoked    bool
    CreatedAt  time.Time
}

type RegisterPayload struct {
    Username	string `validate:"required" json:"username"`
    Email    	string `validate:"required,email" json:"email"`
    Password 	string `validate:"required" json:"password"`
    FirstName 	*string `json:"first_name"`
    LastName 	*string `json:"last_name"`
    Mobile 		*string `json:"mobile_number"`
}

type LoginPayload struct {
    Username    string `validate:"required" json:"username"`
    Password 	string `validate:"required" json:"password"`
}

type RefreshTokenPayload struct {
    RefreshToken string `validate:"required" json:"refresh_token"`
}