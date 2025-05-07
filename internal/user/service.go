package user

import (
	"context"
	"log"
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

func (s *Service) GetUserByID(ctx context.Context, id int) (*UsersList, error) {

	user, err := s.db.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &UsersList{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  user.Fname,
		LastName:   user.Lname,
		Mobile:     user.MobileNumber,
		IsAdmin:    user.IsAdmin,
	}, nil

}


func (s *Service) UpdateUser(ctx context.Context, id int, payload UpdateUser) error {

	if payload.Fname != nil {
		params := db.UpdateUserFirstNameParams{
			ID:        int32(id),
			Fname:     *payload.Fname,
		}
		_, err := s.db.UpdateUserFirstName(ctx, params)
		if err != nil {
			return err
		}
	}

	if payload.Lname != nil {
		params := db.UpdateUserLastNameParams{
			ID:        int32(id),
			Lname:     *payload.Lname,
		}
		_, err := s.db.UpdateUserLastName(ctx, params)
		if err != nil {
			return err
		}
	}

	if payload.MobileNumber != nil {
		params := db.UpdateUserMobileNumberParams{
			ID:        int32(id),
			MobileNumber:     *payload.MobileNumber,
		}
		_, err := s.db.UpdateUserMobileNumber(ctx, params)
		if err != nil {
			return err
		}
	}

	if payload.IsAdmin != nil {
		params := db.UpdateUserRoleParams{
			ID:        int32(id),
			IsAdmin:     *payload.IsAdmin,
		}
		_, err := s.db.UpdateUserRole(ctx, params)
		if err != nil {
			return err
		}
	}

	return nil

}


func (s *Service) DeleteUser(id int) error {

	ID, err := s.db.DeleteUser(context.Background(), int32(id))
	log.Println("err", err)
	log.Println("ID", ID)
	return err

}

func (s *Service) GetUsers(ctx context.Context, limit int, offset int) ([]UsersList, error) {

	params := db.GetUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	
	users, err := s.db.GetUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	var userList []UsersList
	for _, user := range users {
		userList = append(userList, UsersList{
			ID:        	user.ID,
			Username:	user.Username,
			Email:    	user.Email,
			FirstName: 	user.Fname,
			LastName:  	user.Lname,
			Mobile:    	user.MobileNumber,
			IsAdmin:   	user.IsAdmin,
		})
	}

	return userList, nil
}