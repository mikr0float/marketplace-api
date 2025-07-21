package storage

import (
	"context"
	"errors"
	"log"
	"marketAPI/internal/db"
	"marketAPI/internal/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	storage *db.PostgresStorage
}

/*type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	Exists(username string) (bool, error)
}*/

func NewUserRepository(storage *db.PostgresStorage) *UserRepository {
	return &UserRepository{storage: storage}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Bad password: %v", err)
		return err
	}

	user.Password = string(hashedPassword)

	err = sess.InsertInto("users").
		Columns("username", "password").
		Record(user).
		Returning("id", "created_at").
		Load(user)

	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return nil, err
	}

	var user domain.User
	err = sess.Select("*").
		From("users").
		Where("username = ?", username).
		LoadOne(&user)

	if errors.Is(err, dbr.ErrNotFound) {
		return nil, err
	}

	return &user, err
}

func (r *UserRepository) Exists(ctx context.Context, username string) (bool, error) {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return false, err
	}

	var exists bool
	err = sess.Select("1").
		From("users").
		Where("username = ?", username).
		LoadOne(&exists)

	if errors.Is(err, dbr.ErrNotFound) {
		return false, err
	}

	return exists, err
}
