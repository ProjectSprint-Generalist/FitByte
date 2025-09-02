package models

type Activity struct {
	BaseEntity

	ActivityType      string `json:"activity_type"`
	DurationInMinutes int    `json:"duration_in_minutes"`
	CaloriesBurned    int    `json:"calories_burned"`
}
