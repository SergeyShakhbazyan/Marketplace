package repository

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/utils"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id gocql.UUID) (*models.UserWrapContent, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) UserRepository {
	return &userRepository{session: session}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Concatenate countryCode and number into phoneNumber string
	//phoneNumber := user.PhoneNumber.CountryCode + user.PhoneNumber.Number

	query := "INSERT INTO marketplace_keyspace.userData(id, firstName, lastName, email, password, phoneNumber, avatar, AccountType, Subscription, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	return r.session.Query(query,
		user.UserID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.PhoneNumber,
		user.Avatar,
		user.AccountType,
		user.Subscription,
		user.CreatedAt,
	).WithContext(ctx).Exec()
}

func (r *userRepository) GetUser(ctx context.Context, id gocql.UUID) (*models.UserWrapContent, error) {
	//phoneNumber := user.PhoneNumber.CountryCode + user.PhoneNumber.Number
	query := "SELECT id, firstName, lastName, avatar, accountType, rating FROM marketplace_keyspace.userdata WHERE id = ?"
	var user models.UserWrapContent
	err := r.session.Query(query, id).WithContext(ctx).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Avatar, &user.AccountType, &user.Rating)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, firstName, lastName, email, password, phonenumber, avatar, accountType, subscription, createdat FROM marketplace_keyspace.userdata WHERE email = ? ALLOW FILTERING "
	var user models.User
	err := r.session.Query(query, email).WithContext(ctx).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.PhoneNumber, &user.Avatar, &user.AccountType, &user.Subscription, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := "DELETE FROM marketplace_keyspace.userdata WHERE id = ?"
	return r.session.Query(query, id).WithContext(ctx).Exec()
}
