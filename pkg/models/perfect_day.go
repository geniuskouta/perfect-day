package models

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type PerfectDay struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Username    string     `json:"username"`
	Date        string     `json:"date"`
	Areas       []string   `json:"areas"`
	Activities  []Activity `json:"activities"`
	IsDeleted   bool       `json:"is_deleted"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewPerfectDay(id, title, description, username, date string) (*PerfectDay, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	if err := validateDate(date); err != nil {
		return nil, err
	}

	now := time.Now()
	return &PerfectDay{
		ID:          id,
		Title:       title,
		Description: description,
		Username:    username,
		Date:        date,
		Areas:       []string{},
		Activities:  []Activity{},
		IsDeleted:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func validateDate(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %v", err)
	}
	return nil
}

func (pd *PerfectDay) AddActivity(activity Activity) {
	pd.Activities = append(pd.Activities, activity)
	pd.updateAreas()
	pd.UpdatedAt = time.Now()
}

func (pd *PerfectDay) SortActivitiesByTime() {
	sort.Slice(pd.Activities, func(i, j int) bool {
		return pd.Activities[i].StartTime < pd.Activities[j].StartTime
	})
}

func (pd *PerfectDay) updateAreas() {
	areaSet := make(map[string]bool)
	for _, activity := range pd.Activities {
		if activity.Location.Area != "" {
			areaSet[activity.Location.Area] = true
		}
	}

	areas := make([]string, 0, len(areaSet))
	for area := range areaSet {
		areas = append(areas, area)
	}
	sort.Strings(areas)
	pd.Areas = areas
}

func (pd *PerfectDay) UpdateAreas() {
	pd.updateAreas()
	pd.UpdatedAt = time.Now()
}

func (pd *PerfectDay) SoftDelete() {
	pd.IsDeleted = true
	pd.UpdatedAt = time.Now()
}

func (pd *PerfectDay) SearchableContent() string {
	var content strings.Builder
	content.WriteString(pd.Title + " ")
	content.WriteString(pd.Description + " ")

	for _, area := range pd.Areas {
		content.WriteString(area + " ")
	}

	for _, activity := range pd.Activities {
		content.WriteString(activity.Name + " ")
		content.WriteString(activity.Description + " ")
		content.WriteString(activity.Commentary + " ")
		content.WriteString(activity.Location.Name + " ")
	}

	return strings.ToLower(content.String())
}