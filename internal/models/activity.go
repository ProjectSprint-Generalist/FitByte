package models

type Activity struct {
	BaseEntity

	ActivityType      string `gorm:"column:activity_type;type:varchar(100);not null" json:"activity_type"`
	DurationInMinutes int    `gorm:"column:duration_in_minutes;not null" json:"duration_in_minutes"`
	CaloriesBurned    int    `gorm:"column:calories_burned;not null" json:"calories_burned"`
}
