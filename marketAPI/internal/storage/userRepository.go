package storage

import (
	"context"
	"errors"
	"log"
	"marketAPI/internal/db"
	"marketAPI/internal/domain"
	"strings"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	storage *db.PostgresStorage
}

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

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return nil, err
	}

	var user domain.User
	err = sess.Select("*").
		From("users").
		Where("id = ?", id).
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

func (r *UserRepository) PingDB(ctx context.Context) error {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return err
	}
	var result int

	err = sess.Select("1").
		From("users").
		Limit(1).
		LoadOne(&result)
	if err != nil {
		// Если таблицы нет - это тоже "доступная БД"
		if strings.Contains(err.Error(), "relation \"users\" does not exist") {
			return nil
		}
		return err
	}
	return nil
}
