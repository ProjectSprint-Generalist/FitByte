package models

import (
	"time"
)

type Activity struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	ActivityID        string    `json:"activityId" gorm:"uniqueIndex;not null"`
	ActivityType      string    `json:"activityType" gorm:"not null"`
	DoneAt            time.Time `json:"doneAt" gorm:"not null"`
	DurationInMinutes int       `json:"durationInMinutes" gorm:"not null"`
	CaloriesBurned    int       `json:"caloriesBurned" gorm:"not null"`
	UserID            uint      `json:"userId" gorm:"not null;index"`
	User              User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CreateActivityRequest struct {
	ActivityID        string    `json:"activityId" binding:"required"`
	ActivityType      string    `json:"activityType" binding:"required"`
	DoneAt            time.Time `json:"doneAt" binding:"required"`
	DurationInMinutes int       `json:"durationInMinutes" binding:"required,min=1"`
	CaloriesBurned    int       `json:"caloriesBurned" binding:"required,min=0"`
}

type UpdateActivityRequest struct {
	ActivityID        *string    `json:"activityId,omitempty"`
	ActivityType      *string    `json:"activityType,omitempty"`
	DoneAt            *time.Time `json:"doneAt,omitempty"`
	DurationInMinutes *int       `json:"durationInMinutes,omitempty" binding:"omitempty,min=1"`
	CaloriesBurned    *int       `json:"caloriesBurned,omitempty" binding:"omitempty,min=0"`
}

type ActivityResponse struct {
	ID                uint      `json:"id"`
	ActivityID        string    `json:"activityId"`
	ActivityType      string    `json:"activityType"`
	DoneAt            time.Time `json:"doneAt"`
	DurationInMinutes int       `json:"durationInMinutes"`
	CaloriesBurned    int       `json:"caloriesBurned"`
	UserID            uint      `json:"userId"`
	User              UserResponse `json:"user"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
