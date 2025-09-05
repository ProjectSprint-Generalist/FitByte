package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ActivityType string

const (
	ActivityTypeWalking    ActivityType = "Walking"
	ActivityTypeYoga       ActivityType = "Yoga"
	ActivityTypeStretching ActivityType = "Stretching"
	ActivityTypeCycling    ActivityType = "Cycling"
	ActivityTypeSwimming   ActivityType = "Swimming"
	ActivityTypeDancing    ActivityType = "Dancing"
	ActivityTypeHiking     ActivityType = "Hiking"
	ActivityTypeRunning    ActivityType = "Running"
	ActivityTypeHIIT       ActivityType = "HIIT"
	ActivityTypeJumpRope   ActivityType = "JumpRope"
)

type ActivityTypeInfo struct {
	Type            ActivityType `json:"type"`
	CaloriesPerMinute int        `json:"caloriesPerMinute"`
	DisplayName     string       `json:"displayName"`
}

func (at ActivityType) GetCaloriesPerMinute() int {
	switch at {
		case ActivityTypeWalking, ActivityTypeYoga, ActivityTypeStretching:
			return 4
		case ActivityTypeCycling, ActivityTypeSwimming, ActivityTypeDancing:
			return 8
		case ActivityTypeHiking, ActivityTypeRunning, ActivityTypeHIIT, ActivityTypeJumpRope:
			return 10
		default:
			return 0
	}
}

func (at ActivityType) GetDisplayName() string {
	switch at {
		case ActivityTypeWalking:
			return "Walking"
		case ActivityTypeYoga:
			return "Yoga"
		case ActivityTypeStretching:
			return "Stretching"
		case ActivityTypeCycling:
			return "Cycling"
		case ActivityTypeSwimming:
			return "Swimming"
		case ActivityTypeDancing:
			return "Dancing"
		case ActivityTypeHiking:
			return "Hiking"
		case ActivityTypeRunning:
			return "Running"
		case ActivityTypeHIIT:
			return "HIIT"
		case ActivityTypeJumpRope:
			return "Jump Rope"
		default:
			return string(at)
	}
}

func (at ActivityType) IsValid() bool {
	switch at {
		case ActivityTypeWalking, ActivityTypeYoga, ActivityTypeStretching,
			ActivityTypeCycling, ActivityTypeSwimming, ActivityTypeDancing,
			ActivityTypeHiking, ActivityTypeRunning, ActivityTypeHIIT, ActivityTypeJumpRope:
			return true
		default:
			return false
	}
}

func ParseActivityType(s string) (ActivityType, error) {
	at := ActivityType(strings.TrimSpace(s))
	if !at.IsValid() {
		return "", fmt.Errorf("invalid activity type: %s", s)
	}
	return at, nil
}

func GetAllActivityTypes() []ActivityTypeInfo {
	return []ActivityTypeInfo{
		{Type: ActivityTypeWalking, CaloriesPerMinute: ActivityTypeWalking.GetCaloriesPerMinute(), DisplayName: ActivityTypeWalking.GetDisplayName()},
		{Type: ActivityTypeYoga, CaloriesPerMinute: ActivityTypeYoga.GetCaloriesPerMinute(), DisplayName: ActivityTypeYoga.GetDisplayName()},
		{Type: ActivityTypeStretching, CaloriesPerMinute: ActivityTypeStretching.GetCaloriesPerMinute(), DisplayName: ActivityTypeStretching.GetDisplayName()},
		{Type: ActivityTypeCycling, CaloriesPerMinute: ActivityTypeCycling.GetCaloriesPerMinute(), DisplayName: ActivityTypeCycling.GetDisplayName()},
		{Type: ActivityTypeSwimming, CaloriesPerMinute: ActivityTypeSwimming.GetCaloriesPerMinute(), DisplayName: ActivityTypeSwimming.GetDisplayName()},
		{Type: ActivityTypeDancing, CaloriesPerMinute: ActivityTypeDancing.GetCaloriesPerMinute(), DisplayName: ActivityTypeDancing.GetDisplayName()},
		{Type: ActivityTypeHiking, CaloriesPerMinute: ActivityTypeHiking.GetCaloriesPerMinute(), DisplayName: ActivityTypeHiking.GetDisplayName()},
		{Type: ActivityTypeRunning, CaloriesPerMinute: ActivityTypeRunning.GetCaloriesPerMinute(), DisplayName: ActivityTypeRunning.GetDisplayName()},
		{Type: ActivityTypeHIIT, CaloriesPerMinute: ActivityTypeHIIT.GetCaloriesPerMinute(), DisplayName: ActivityTypeHIIT.GetDisplayName()},
		{Type: ActivityTypeJumpRope, CaloriesPerMinute: ActivityTypeJumpRope.GetCaloriesPerMinute(), DisplayName: ActivityTypeJumpRope.GetDisplayName()},
	}
}

// GenerateActivityID generates a unique activity ID based on the activity type
func GenerateActivityID(activityType ActivityType) string {
	// Get a short prefix based on activity type
	var prefix string
	switch activityType {
	case ActivityTypeWalking:
		prefix = "WALK"
	case ActivityTypeYoga:
		prefix = "YOGA"
	case ActivityTypeStretching:
		prefix = "STRETCH"
	case ActivityTypeCycling:
		prefix = "CYCLE"
	case ActivityTypeSwimming:
		prefix = "SWIM"
	case ActivityTypeDancing:
		prefix = "DANCE"
	case ActivityTypeHiking:
		prefix = "HIKE"
	case ActivityTypeRunning:
		prefix = "RUN"
	case ActivityTypeHIIT:
		prefix = "HIIT"
	case ActivityTypeJumpRope:
		prefix = "JUMP"
	default:
		prefix = "ACT"
	}
	
	// Generate a random 8-character string
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	
	randomPart := make([]byte, 8)
	for i := range randomPart {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}
	
	// Add timestamp for additional uniqueness
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%s_%s_%d", prefix, string(randomPart), timestamp)
}

type Activity struct {
	ID                uint        `json:"id" gorm:"primaryKey"`
	ActivityID        string      `json:"activityId" gorm:"uniqueIndex;not null"`
	ActivityType      ActivityType `json:"activityType" gorm:"not null"`
	DoneAt            time.Time   `json:"doneAt" gorm:"not null"`
	DurationInMinutes int         `json:"durationInMinutes" gorm:"not null"`
	CaloriesBurned    int         `json:"caloriesBurned" gorm:"not null"`
	UserID            uint        `json:"userId" gorm:"not null;index"`
	User              User        `json:"user" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

type CreateActivityRequest struct {
	ActivityID        *string     `json:"activityId,omitempty"`
	ActivityType      string      `json:"activityType" binding:"required"`
	DoneAt            time.Time   `json:"doneAt" binding:"required"`
	DurationInMinutes int         `json:"durationInMinutes" binding:"required,min=1"`
	CaloriesBurned    *int        `json:"caloriesBurned,omitempty" binding:"omitempty,min=0"`
}

func (req *CreateActivityRequest) Validate() error {
	activityType, err := ParseActivityType(req.ActivityType)
	if err != nil {
		return err
	}
	
	// Generate activity ID if not provided
	if req.ActivityID == nil {
		activityID := GenerateActivityID(activityType)
		req.ActivityID = &activityID
	}
	
	// Calculate calories if not provided
	if req.CaloriesBurned == nil {
		calories := activityType.GetCaloriesPerMinute() * req.DurationInMinutes
		req.CaloriesBurned = &calories
	}
	
	return nil
}

type UpdateActivityRequest struct {
	ActivityID        *string    `json:"activityId,omitempty"`
	ActivityType      *string    `json:"activityType,omitempty"`
	DoneAt            *time.Time `json:"doneAt,omitempty"`
	DurationInMinutes *int       `json:"durationInMinutes,omitempty" binding:"omitempty,min=1"`
	CaloriesBurned    *int       `json:"caloriesBurned,omitempty" binding:"omitempty,min=0"`
}

func (req *UpdateActivityRequest) Validate() error {
	// If activity type is being updated, validate it
	if req.ActivityType != nil {
		_, err := ParseActivityType(*req.ActivityType)
		if err != nil {
			return err
		}
	}
	
	return nil
}

type ActivityResponse struct {
	ID                uint        `json:"id"`
	ActivityID        string      `json:"activityId"`
	ActivityType      ActivityType `json:"activityType"`
	DoneAt            time.Time   `json:"doneAt"`
	DurationInMinutes int         `json:"durationInMinutes"`
	CaloriesBurned    int         `json:"caloriesBurned"`
	UserID            uint        `json:"userId"`
	User              UserResponse `json:"user"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}
