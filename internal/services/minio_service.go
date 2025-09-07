package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"fitbyte/internal/config"
	"fitbyte/internal/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type MinIOService struct {
	client     *minio.Client
	bucketName string
	db         *gorm.DB
}

func NewMinIOService(cfg *config.Config, db *gorm.DB) (*MinIOService, error) {
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	service := &MinIOService{
		client:     client,
		bucketName: cfg.MinIO.BucketName,
		db:         db,
	}

	if err := service.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return service, nil
}

func (s *MinIOService) ensureBucketExists() error {
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MinIOService) UploadFile(file *multipart.FileHeader, userID uint) (*models.FileUploadResponse, error) {
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	objectName := fmt.Sprintf("%d/%s_%d%s", userID, strings.TrimSuffix(file.Filename, fileExt), time.Now().Unix(), fileExt)

	ctx := context.Background()
	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	fileURI := fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), s.bucketName, objectName)

	fileUpload := &models.FileUpload{
		FileName: file.Filename,
		FileSize: file.Size,
		FileType: file.Header.Get("Content-Type"),
		FileURI:  fileURI,
		UserID:   userID,
	}

	if err := s.db.Create(fileUpload).Error; err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return &models.FileUploadResponse{
		URI: fileURI,
	}, nil
}

func (s *MinIOService) validateFile(file *multipart.FileHeader) error {
	const maxSize = 100 * 1024 // 100 KiB
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds 100 KiB limit")
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return fmt.Errorf("file type %s not allowed. Only JPEG, JPG, and PNG are supported", contentType)
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range allowedExts {
		if fileExt == ext {
			return nil
		}
	}

	return fmt.Errorf("file extension %s not allowed. Only .jpg, .jpeg, and .png are supported", fileExt)
}

func (s *MinIOService) DeleteFile(fileURI string, userID uint) error {
	fileUpload := &models.FileUpload{}
	if err := s.db.Where("file_uri = ? AND user_id = ?", fileURI, userID).First(fileUpload).Error; err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	objectName := strings.TrimPrefix(fileUpload.FileURI, fmt.Sprintf("%s/%s/", s.client.EndpointURL(), s.bucketName))
	ctx := context.Background()
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	if err := s.db.Delete(fileUpload).Error; err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}

	return nil
}

func (s *MinIOService) GetUserFiles(userID uint, limit, offset int) ([]models.FileUpload, error) {
	var files []models.FileUpload
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error
	return files, err
}
