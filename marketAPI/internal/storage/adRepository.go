package storage

import (
	"context"
	"log"
	"marketAPI/internal/db"
	"marketAPI/internal/domain"

	"github.com/gocraft/dbr/v2"
)

/*type AdRepository interface {
	Create(ad *domain.Ad) error
	GetByID(id int) (*domain.Ad, error)
	List(filter domain.AdFilter) ([]domain.Ad, error)
}*/

type AdRepository struct {
	storage *db.PostgresStorage
}

func NewAdRepository(storage *db.PostgresStorage) *AdRepository {
	return &AdRepository{storage: storage}
}

func (r *AdRepository) Create(ctx context.Context, ad *domain.Ad) error {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return err
	}

	err = sess.InsertInto("ads").
		Columns("title", "description", "image_url", "price", "user_id").
		Record(ad).
		Returning("id", "created_at").
		Load(ad)

	return err
}

func (r *AdRepository) GetByID(ctx context.Context, id int) (*domain.Ad, error) {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return nil, err
	}

	var ad domain.Ad
	err = sess.Select("*").
		From("ads").
		Where("id = ?", id).
		LoadOne(&ad)

	return &ad, err
}

func (r *AdRepository) List(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
	sess, err := r.storage.NewSession(ctx)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return nil, err
	}

	query := sess.Select("a.*", "u.username").
		From(dbr.I("ads").As("a")).
		Join(dbr.I("users").As("u"), "a.user_id = u.id")

	if filter.MinPrice != nil {
		query.Where("a.price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query.Where("a.price <= ?", *filter.MaxPrice)
	}

	// Сортировка
	sortField := "a.created_at"
	if filter.SortBy == "price" {
		sortField = "a.price"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query.OrderDir(sortField, sortOrder == "ASC")

	// Пагинация
	if filter.PageSize > 0 {
		query.Limit(uint64(filter.PageSize))
		if filter.Page > 1 {
			query.Offset(uint64((filter.Page - 1) * filter.PageSize))
		}
	}

	var ads []domain.Ad
	_, err = query.Load(&ads)

	return ads, err
}
