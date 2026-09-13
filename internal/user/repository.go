package user

import (
	"errors"
	"student-go-service/internal/package/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
func (r *repository) Create(user *User) error {
	err := r.db.Create(user).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.ErrEmailAlreadyExists
		}

		return err
	}

	return nil
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User

	if err := r.db.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
