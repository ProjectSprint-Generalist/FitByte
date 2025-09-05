package models

// ActivityType represents the enum for activity types
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

// CreateActivityRequest represents the request payload for creating an activity
type CreateActivityRequest struct {
	ActivityType      ActivityType `json:"activityType" binding:"required"`
	DoneAt            time.Time    `json:"doneAt" binding:"required"`
	DurationInMinutes int          `json:"durationInMinutes" binding:"required,min=1"`
}

// ActivityResponse represents the response payload for activity data
type ActivityResponse struct {
	ActivityId        uint         `json:"activityId"`
	ActivityType      ActivityType `json:"activityType"`
	DoneAt            time.Time    `json:"doneAt"`
	DurationInMinutes int          `json:"durationInMinutes"`
	CaloriesBurned    int          `json:"caloriesBurned"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

// CaloriesPerMinute defines the calories burned per minute for each activity type
var CaloriesPerMinute = map[ActivityType]int{
	ActivityTypeWalking:    4,
	ActivityTypeYoga:       4,
	ActivityTypeStretching: 4,
	ActivityTypeCycling:    8,
	ActivityTypeSwimming:   8,
	ActivityTypeDancing:    8,
	ActivityTypeHiking:     10,
	ActivityTypeRunning:    10,
	ActivityTypeHIIT:       10,
	ActivityTypeJumpRope:   10,
}

// CalculateCalories calculates calories burned based on activity type and duration
func (a *Activity) CalculateCalories() int {
	caloriesPerMin, exists := CaloriesPerMinute[a.ActivityType]
	if !exists {
		return 0
	}
	return caloriesPerMin * a.DurationInMinutes
}

// ToResponse converts Activity to ActivityResponse
func (a *Activity) ToResponse() ActivityResponse {
	return ActivityResponse{
		ActivityId:        a.ID,
		ActivityType:      a.ActivityType,
		DoneAt:            a.DoneAt,
		DurationInMinutes: a.DurationInMinutes,
		CaloriesBurned:    a.CaloriesBurned,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

// IsValidActivityType checks if the activity type is valid
func IsValidActivityType(activityType string) bool {
	validTypes := []ActivityType{
		ActivityTypeWalking,
		ActivityTypeYoga,
		ActivityTypeStretching,
		ActivityTypeCycling,
		ActivityTypeSwimming,
		ActivityTypeDancing,
		ActivityTypeHiking,
		ActivityTypeRunning,
		ActivityTypeHIIT,
		ActivityTypeJumpRope,
	}

	for _, validType := range validTypes {
		if string(validType) == activityType {
			return true
		}
	}
	return false
}

type Activity struct {
	BaseEntity

	ActivityType      string    `gorm:"column:activity_type;type:varchar(100);not null" json:"activity_type"`
	DurationInMinutes int       `gorm:"column:duration_in_minutes;not null" json:"duration_in_minutes"`
	CaloriesBurned    int       `gorm:"column:calories_burned;not null" json:"calories_burned"`
	DoneAt            time.Time `gorm:"column:done_at;not null" json:"done_at"`
}

type UpdateActivityRequest struct {
	ActivityType      *string `gorm:"column:activity_type;type:varchar(100);not null" json:"activity_type"`
	DurationInMinutes *int    `gorm:"column:duration_in_minutes;not null" json:"duration_in_minutes"`
	CaloriesBurned    *int    `gorm:"column:calories_burned;not null" json:"calories_burned"`
}
