package models

import "time"

type FileUpload struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FileName  string    `json:"file_name" gorm:"not null"`
	FileSize  int64     `json:"file_size" gorm:"not null"`
	FileType  string    `json:"file_type" gorm:"not null"`
	FileURI   string    `json:"file_uri" gorm:"not null;unique"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileUploadRequest struct {
	File string `form:"file" binding:"required"`
}

type FileUploadResponse struct {
	URI string `json:"uri"`
}
