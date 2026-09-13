package utils

import (
	"fmt"
	"mime/multipart"
	"strings"
)

const MaxImageSize = 5 * 1024 * 1024 //5mb

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func ValidateImage(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return nil
	}

	if fileHeader.Size > MaxImageSize {
		return fmt.Errorf("Image Size Must not exceed 5 MB")
	}
	contentType := strings.ToLower(fileHeader.Header.Get("Content-Type"))

	if !allowedImageTypes[contentType] {
		return fmt.Errorf(
			"invalid image type :%s; allowed types are JPEG , PNG and WEBP", contentType)
	}
	return nil
}
