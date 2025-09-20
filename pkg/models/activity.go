package models

import (
	"fmt"
	"time"
)

type Activity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Location    Location  `json:"location"`
	StartTime   string    `json:"start_time"`
	Duration    int       `json:"duration_minutes"`
	Description string    `json:"description,omitempty"`
	Commentary  string    `json:"commentary,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewActivity(id, name string, location Location, startTime string, duration int, description, commentary string) (*Activity, error) {
	if err := validateActivityTime(startTime); err != nil {
		return nil, err
	}

	if duration <= 0 {
		return nil, fmt.Errorf("duration must be positive")
	}

	if name == "" {
		return nil, fmt.Errorf("activity name is required")
	}

	return &Activity{
		ID:          id,
		Name:        name,
		Location:    location,
		StartTime:   startTime,
		Duration:    duration,
		Description: description,
		Commentary:  commentary,
		CreatedAt:   time.Now(),
	}, nil
}

func validateActivityTime(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:MM: %v", err)
	}
	return nil
}

func (a *Activity) EndTime() (string, error) {
	startTime, err := time.Parse("15:04", a.StartTime)
	if err != nil {
		return "", err
	}

	endTime := startTime.Add(time.Duration(a.Duration) * time.Minute)
	return endTime.Format("15:04"), nil
}