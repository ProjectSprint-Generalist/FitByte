package services

import (
	"errors"
	"fmt"

	"fitbyte/internal/database"
	"fitbyte/internal/models"

	"gorm.io/gorm"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService() *ActivityService {
	return &ActivityService{
		db: database.DB,
	}
}

func (s *ActivityService) CreateActivity(userID uint, req *models.CreateActivityRequest) (*models.Activity, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	activityType, err := models.ParseActivityType(req.ActivityType)
	if err != nil {
		return nil, fmt.Errorf("invalid activity type: %w", err)
	}

	caloriesBurned := *req.CaloriesBurned
	if caloriesBurned == 0 {
		caloriesBurned = activityType.GetCaloriesPerMinute() * req.DurationInMinutes
	}

	activity := &models.Activity{
		ActivityID:        *req.ActivityID,
		ActivityType:      activityType,
		DoneAt:            req.DoneAt,
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		UserID:            userID,
	}

	if err := s.db.Create(activity).Error; err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return activity, nil
}

func (s *ActivityService) GetActivityByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	if err := s.db.Preload("User").First(&activity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("activity not found")
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	return &activity, nil
}

func (s *ActivityService) GetActivityByActivityID(activityID string) (*models.Activity, error) {
	var activity models.Activity
	if err := s.db.Preload("User").Where("activity_id = ?", activityID).First(&activity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("activity not found")
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	return &activity, nil
}

func (s *ActivityService) GetUserActivities(userID uint, page, limit int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	query := s.db.Model(&models.Activity{}).Where("user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Offset(offset).Limit(limit).Order("done_at DESC").Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get activities: %w", err)
	}

	return activities, total, nil
}

func (s *ActivityService) GetAllActivities(page, limit int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	if err := s.db.Model(&models.Activity{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	offset := (page - 1) * limit
	if err := s.db.Preload("User").Offset(offset).Limit(limit).Order("done_at DESC").Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get activities: %w", err)
	}

	return activities, total, nil
}

func (s *ActivityService) UpdateActivity(id uint, req *models.UpdateActivityRequest) (*models.Activity, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	var activity models.Activity
	if err := s.db.Preload("User").First(&activity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("activity not found")
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	if req.ActivityID != nil {
		activity.ActivityID = *req.ActivityID
	}
	if req.ActivityType != nil {
		activityType, err := models.ParseActivityType(*req.ActivityType)
		if err != nil {
			return nil, fmt.Errorf("invalid activity type: %w", err)
		}
		activity.ActivityType = activityType
	}
	if req.DoneAt != nil {
		activity.DoneAt = *req.DoneAt
	}
	if req.DurationInMinutes != nil {
		activity.DurationInMinutes = *req.DurationInMinutes
	}
	if req.CaloriesBurned != nil {
		activity.CaloriesBurned = *req.CaloriesBurned
	}

	if (req.ActivityType != nil || req.DurationInMinutes != nil) && req.CaloriesBurned == nil {
		activity.CaloriesBurned = activity.ActivityType.GetCaloriesPerMinute() * activity.DurationInMinutes
	}

	if err := s.db.Save(&activity).Error; err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	return &activity, nil
}

func (s *ActivityService) DeleteActivity(id uint) error {
	if err := s.db.Delete(&models.Activity{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	return nil
}

func (s *ActivityService) GetActivitiesByType(activityType string, page, limit int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	query := s.db.Model(&models.Activity{}).Where("activity_type = ?", activityType)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Offset(offset).Limit(limit).Order("done_at DESC").Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get activities: %w", err)
	}

	return activities, total, nil
}

func (s *ActivityService) GetUserActivitiesByType(userID uint, activityType string, page, limit int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	query := s.db.Model(&models.Activity{}).Where("user_id = ? AND activity_type = ?", userID, activityType)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Offset(offset).Limit(limit).Order("done_at DESC").Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get activities: %w", err)
	}

	return activities, total, nil
}

func (s *ActivityService) ToActivityResponse(activity *models.Activity) models.ActivityResponse {
	userService := NewUserService()
	return models.ActivityResponse{
		ID:                activity.ID,
		ActivityID:        activity.ActivityID,
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt,
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		UserID:            activity.UserID,
		User:              userService.ToUserResponse(&activity.User),
		CreatedAt:         activity.CreatedAt,
		UpdatedAt:         activity.UpdatedAt,
	}
}

func (s *ActivityService) GetActivityTypes() []models.ActivityTypeInfo {
	return models.GetAllActivityTypes()
}
