package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage keeps the S3 client and the bucket name.
type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
}

// NewS3Storage creates a new S3 storage object.
func NewS3Storage(
	cfg aws.Config,
	bucketName string,
) *S3Storage {
	client := s3.NewFromConfig(cfg)
	return &S3Storage{
		// Make an S3 client using the AWS settings.
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    bucketName,
	}

}

// Upload sends a file to the S3 bucket.
//
// key is the file name or path inside the bucket.
func (s *S3Storage) Upload(
	ctx context.Context,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	key string,
) error {
	// PutObject sends the file to S3.
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		// Choose the bucket.
		Bucket: aws.String(s.bucketName),
		// Choose the file name or path.
		Key: aws.String(key),
		// Give S3 the file data.
		Body: file,
		// Tell S3 what type of file this is.
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file to s3: %w", err)
	}

	// Return an error if the upload failed.
	return nil
}
func GetFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func (s *S3Storage) GetPresignedURL(
	ctx context.Context,
	key string,
	expires time.Duration,
) (string, error) {
	result, err := s.presignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = expires
		},
	)
	if err != nil {
		return "", fmt.Errorf("Failed to Generate Presigned URL %w", err)
	}
	return result.URL, nil
}

func (s *S3Storage) TestConnection(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return fmt.Errorf("Failed to connect to S3 Bucket: %w", err)
	}
	return nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}

	_, err := s.client.DeleteObject(
		ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete file from s3:%w", err)
	}
	return nil

}
