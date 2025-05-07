package auth

import (
	"context"
	"time"

	db "user-management/internal/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)


type Service struct {
	db *db.Queries
}

func NewService(db *db.Queries) *Service {
	return &Service{
		db: db,
	}
}

func ToPgText(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func (s *Service) Register(ctx context.Context, payload struct {
	Username	string `validate:"required" json:"username"`
	Email    	string `validate:"required,email" json:"email"`
	Password 	string `validate:"required" json:"password"`
	FirstName 	*string `json:"first_name"`
	LastName 	*string `json:"last_name"`
	Mobile 		*string `json:"mobile_number"`
}) error {

	params := db.CreateUserParams{
		Username:    	payload.Username,
		Email: 			payload.Email,
		PasswordHash: 	payload.Password,
		Fname: 			ToPgText(payload.FirstName).String,
		Lname: 			ToPgText(payload.LastName).String,
		MobileNumber: 	ToPgText(payload.Mobile).String,
	}

	_, err := s.db.CreateUser(ctx, params)
	return err

}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {

	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (s *Service) GetUserByID(ctx context.Context, id int) (*db.User, error) {

	user, err := s.db.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (s *Service) SaveRefreshToken(ctx context.Context, payload struct {
	UserID		int32 		`json:"user_id"`
	TokenHash   string 		`json:"token_hash"`
	ExpiresAt 	time.Time 	`json:"expires_at"`
}) error {

	params := db.CreateRefreshTokenParams{
		UserID:    	payload.UserID,
		TokenHash: 	payload.TokenHash,
		ExpiresAt: 	payload.ExpiresAt,
	}

	_, err := s.db.CreateRefreshToken(context.Background(), params)
	return err

}

func (s *Service) FindRefreshToken(ctx context.Context, tokenHash string) (*db.RefreshToken, error) {

	refresh, err := s.db.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &refresh, nil

}

func (s *Service) DeleteRefreshTokenByUserID(ctx context.Context, userID int32) error {

	err := s.db.RevokeRefreshToken(ctx, int32(userID))
	return err

}