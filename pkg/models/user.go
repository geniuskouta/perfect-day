package models

import (
	"fmt"
	"regexp"
	"time"
)

type User struct {
	Username  string    `json:"username"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(username, timezone string) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validateTimezone(timezone); err != nil {
		return nil, err
	}

	return &User{
		Username:  username,
		Timezone:  timezone,
		CreatedAt: time.Now(),
	}, nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be 3-20 characters long")
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, hyphens, and underscores")
	}

	return nil
}

func validateTimezone(timezone string) error {
	_, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %v", err)
	}
	return nil
}