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