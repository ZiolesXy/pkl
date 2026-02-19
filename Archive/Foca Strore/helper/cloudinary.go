package helper

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloudinaryInstance *cloudinary.Cloudinary

// Struct untuk menampung hasil upload agar bisa diakses di handler
type UploadResult struct {
	SecureURL string
	PublicID  string
}

// Helper untuk menangani *bool
func ptrBool(b bool) *bool {
	return &b
}

func InitCloudinary() error {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return errors.New("cloudinary credentials not found")
	}

	var err error
	cloudinaryInstance, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return nil
}

func UploadImageFromFile(filePath string, folder string) (*UploadResult, error) {
	if cloudinaryInstance == nil {
		return nil, errors.New("cloudinary not initialized")
	}

	uploadParams := uploader.UploadParams{
		Folder:         folder,
		UseFilename:    ptrBool(true),
		UniqueFilename: ptrBool(true),
		Overwrite:      ptrBool(false),
	}

	resp, err := cloudinaryInstance.Upload.Upload(context.Background(), filePath, uploadParams)
	if err != nil {
		return nil, err
	}

	return &UploadResult{SecureURL: resp.SecureURL, PublicID: resp.PublicID}, nil
}

func UploadImageFromURL(imageURL string, folder string) (*UploadResult, error) {
	if cloudinaryInstance == nil {
		return nil, errors.New("cloudinary not initialized")
	}

	uploadParams := uploader.UploadParams{
		Folder:         folder,
		UseFilename:    ptrBool(true),
		UniqueFilename: ptrBool(true),
		Overwrite:      ptrBool(false),
	}

	resp, err := cloudinaryInstance.Upload.Upload(context.Background(), imageURL, uploadParams)
	if err != nil {
		return nil, err
	}

	return &UploadResult{SecureURL: resp.SecureURL, PublicID: resp.PublicID}, nil
}

func DeleteImage(publicID string) error {
	if cloudinaryInstance == nil {
		return errors.New("cloudinary not initialized")
	}

	// Perbaikan: Gunakan struct DestroyParams dan masukkan PublicID di dalamnya
	_, err := cloudinaryInstance.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}