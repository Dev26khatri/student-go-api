package student

import (
	"errors"
	"log"
	"student-go-service/internal/package/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines database operations for student records.
type Repository interface {
	Create(student *Student) error
	GetAll(pagination utils.Pagination, search, sortBy, order string) ([]Student, int64, error)
	GetById(id uint) (*Student, error)
	Update(student *Student) error
	Delete(id uint) error
}

// repository provides the database implementation of Repository.
type repository struct {
	db *gorm.DB
}

// NewRepository creates a repository backed by the provided database.
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

var ErrEmailAlreadyExist = errors.New("Email Already Is Exists")

// Create inserts a new student record.
func (r *repository) Create(student *Student) error {
	err := r.db.Create(student).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEmailAlreadyExist
			}

		}
		return err

	}
	return nil
}

// GetAll searches for students and returns sorted, paginated records with their total.
func (r *repository) GetAll(pagination utils.Pagination, search, sortBy, order string) ([]Student, int64, error) {
	var student []Student
	var total int64

	query := r.db.Model(&Student{})
	if search != "" {
		searchPattern := "%" + search + "%"

		query = query.Where("name ILIKE ? or email ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order(sortBy + " " + order).
		Find(&student).Error
	return student, total, err
}

// GetById loads one student by its database ID.
func (r *repository) GetById(id uint) (*Student, error) {
	var student Student

	err := r.db.Where("id = ?", id).First(&student).Error

	log.Println("Repo", err)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// Update saves changes to an existing student record.
func (r *repository) Update(student *Student) error {
	err := r.db.Save(&student).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEmailAlreadyExist
			}

		}
		return err
	}
	return nil
}

// Delete removes a student and reports when no record was found.
func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&Student{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
