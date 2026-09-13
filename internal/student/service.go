package student

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"student-go-service/internal/dto"
	"student-go-service/internal/package/utils"
	"student-go-service/internal/storage"
	"time"

	"github.com/google/uuid"
)

// Service connects student business logic to the repository.
type Service struct {
	repository Repository
	storage    *storage.S3Storage
}

// NewService creates a service with the provided repository.
func NewService(repository Repository, storage *storage.S3Storage) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
	}
}

// Create builds and saves a new student.
func (s *Service) Create(req dto.CreateStudentRequest) (*Student, error) {
	student := &Student{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	if err := s.repository.Create(student); err != nil {
		return nil, err
	}

	if req.Photo != nil {

		file, err := req.Photo.Open()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to open uploaded photo: %w",
				err,
			)
		}
		defer file.Close()

		extension := filepath.Ext(req.Photo.Filename)

		key := fmt.Sprintf(
			"students/%d/%s%s",
			student.ID,
			uuid.New().String(),
			extension,
		)

		// Upload to S3
		if err := s.storage.Upload(
			context.Background(),
			file,
			req.Photo,
			key,
		); err != nil {
			return nil, err
		}

		// Save S3 key in student
		student.Photo = key

		// Update PostgreSQL
		if err := s.repository.Update(student); err != nil {

			// DB failed, so remove the S3 object we just created.
			if deleteErr := s.storage.Delete(
				context.Background(),
				key,
			); deleteErr != nil {
				log.Printf(
					"WARNING: failed to cleanup S3 object %s: %v",
					key,
					deleteErr,
				)
			}

			return nil, err
		}
	}

	log.Printf("Student Created Successfully: %d", student.ID)
	return student, nil
}

// GetAll retrieves students with pagination, search, and sorting.
func (s *Service) GetAll(ctx context.Context, pagination utils.Pagination, search, sortBy, order string) ([]dto.StudentResponse, int64, error) {
	students, total, err := s.repository.GetAll(pagination, search, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.StudentResponse, 0, len(students))
	for i := range students {
		item, err := s.toStudentResponse(ctx, &students[i])
		if err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}
	return result, total, nil
}

// GetById retrieves one student by ID.
func (s *Service) GetById(id uint) (*Student, error) {
	log.Println("Service", id)
	return s.repository.GetById(id)
}

// Update changes and saves an existing student.
func (s *Service) Update(
	id uint,
	req *dto.UpdateStudentRequest,
) (*Student, error) {

	// 1. Get existing student
	student, err := s.repository.GetById(id)
	if err != nil {
		return nil, err
	}

	// Keep the old S3 key.
	oldPhoto := student.Photo

	// 2. Update normal fields
	student.Name = req.Name
	student.Email = req.Email
	student.Age = req.Age

	var newPhotoKey string

	// 3. Upload new photo if provided
	if req.Photo != nil {

		file, err := req.Photo.Open()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to open uploaded photo: %w",
				err,
			)
		}
		defer file.Close()

		extension := filepath.Ext(req.Photo.Filename)

		newPhotoKey = fmt.Sprintf(
			"students/%d/%s%s",
			student.ID,
			uuid.New().String(),
			extension,
		)

		// Upload NEW image.
		if err := s.storage.Upload(
			context.Background(),
			file,
			req.Photo,
			newPhotoKey,
		); err != nil {
			return nil, err
		}

		// Change DB value to NEW photo.
		student.Photo = newPhotoKey
	}

	// 4. Save student to PostgreSQL
	if err := s.repository.Update(student); err != nil {

		// DB failed after S3 upload.
		// Delete newly uploaded image to avoid orphan S3 object.
		if newPhotoKey != "" {
			_ = s.storage.Delete(
				context.Background(),
				newPhotoKey,
			)
		}

		return nil, err
	}

	// 5. DB successfully updated.
	// Now delete OLD image.
	if newPhotoKey != "" && oldPhoto != "" {

		if err := s.storage.Delete(
			context.Background(),
			oldPhoto,
		); err != nil {

			// Student update succeeded.
			// Only old S3 cleanup failed.
			log.Printf(
				"WARNING: failed to delete old S3 photo %s: %v",
				oldPhoto,
				err,
			)
		}
	}

	return student, nil
}

// Delete removes a student by ID.
func (s *Service) Delete(id uint) error {
	student, err := s.repository.GetById(id)
	if err != nil {
		return err
	}
	if student.Photo != "" {
		if err := s.storage.Delete(
			context.Background(),
			student.Photo,
		); err != nil {
			return err
		}
	}
	return s.repository.Delete(id)
}

func (s *Service) GetAllPhotoUrl(ctx context.Context, student *Student) (string, error) {
	if student.Photo == "" {
		return "", nil
	}
	return s.storage.GetPresignedURL(ctx, student.Photo, 15*time.Minute)
}

func (s *Service) toStudentResponse(
	ctx context.Context,
	student *Student,
) (dto.StudentResponse, error) {

	photoURL, err := s.GetAllPhotoUrl(ctx, student)
	if err != nil {
		return dto.StudentResponse{}, err
	}

	return dto.StudentResponse{
		ID:        student.ID,
		Name:      student.Name,
		Email:     student.Email,
		Age:       student.Age,
		PhotoUrl:  photoURL,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
	}, nil
}

func (s *Service) GetPhotoUrl(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", nil
	}

	return s.storage.GetPresignedURL(
		ctx, key, 15*time.Minute,
	)
}
