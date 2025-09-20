package unit

import (
	"perfect-day/pkg/models"
	"perfect-day/pkg/search"
	"testing"
)

func createTestPerfectDays() []*models.PerfectDay {
	pd1, _ := models.NewPerfectDay("id1", "Coffee Day in Tokyo", "Great coffee exploration", "alice", "2023-12-01")
	location1 := models.NewCustomTextLocation("Blue Bottle Coffee", "Shibuya")
	activity1, _ := models.NewActivity("act1", "Coffee tasting", *location1, "10:00", 60, "Trying new beans", "Amazing experience")
	pd1.AddActivity(*activity1)

	pd2, _ := models.NewPerfectDay("id2", "Food Tour", "Exploring local cuisine", "bob", "2023-12-02")
	location2 := models.NewCustomTextLocation("Sushi Restaurant", "Shinjuku")
	activity2, _ := models.NewActivity("act2", "Sushi dinner", *location2, "19:00", 120, "Fresh sashimi", "Best sushi ever")
	pd2.AddActivity(*activity2)

	pd3, _ := models.NewPerfectDay("id3", "Museum Visit", "Art and culture day", "alice", "2023-11-30")
	location3 := models.NewCustomTextLocation("National Museum", "Ueno")
	activity3, _ := models.NewActivity("act3", "Art viewing", *location3, "14:00", 180, "Contemporary art exhibit", "Very inspiring")
	pd3.AddActivity(*activity3)

	pd4, _ := models.NewPerfectDay("id4", "Coffee and Books", "Quiet reading day", "charlie", "2023-12-03")
	location4 := models.NewCustomTextLocation("Bookstore Cafe", "Shibuya")
	activity4, _ := models.NewActivity("act4", "Reading", *location4, "13:00", 150, "Philosophy book", "Deep thoughts")
	pd4.AddActivity(*activity4)
	pd4.SoftDelete()

	return []*models.PerfectDay{pd1, pd2, pd3, pd4}
}

func TestSearchByQuery(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{"coffee search", "coffee", 2},
		{"shibuya search", "shibuya", 2},
		{"sushi search", "sushi", 1},
		{"art search", "art", 1},
		{"tokyo search", "tokyo", 1},
		{"nonexistent search", "pizza", 0},
		{"empty query", "", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := search.SearchCriteria{Query: tt.query}
			results := searchService.Search(perfectDays, criteria)

			if results.Total != tt.expectedCount {
				t.Errorf("Expected %d results for query '%s', got %d", tt.expectedCount, tt.query, results.Total)
			}
		})
	}
}

func TestSearchByUser(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	tests := []struct {
		name          string
		username      string
		expectedCount int
	}{
		{"alice search", "alice", 2},
		{"bob search", "bob", 1},
		{"charlie search", "charlie", 1},
		{"nonexistent user", "david", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := search.SearchCriteria{Username: tt.username}
			results := searchService.Search(perfectDays, criteria)

			if results.Total != tt.expectedCount {
				t.Errorf("Expected %d results for user '%s', got %d", tt.expectedCount, tt.username, results.Total)
			}
		})
	}
}

func TestSearchByArea(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	tests := []struct {
		name          string
		areas         []string
		expectedCount int
	}{
		{"shibuya area", []string{"Shibuya"}, 2},
		{"shinjuku area", []string{"Shinjuku"}, 1},
		{"ueno area", []string{"Ueno"}, 1},
		{"multiple areas", []string{"Shibuya", "Ueno"}, 3},
		{"nonexistent area", []string{"Harajuku"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := search.SearchCriteria{Areas: tt.areas}
			results := searchService.Search(perfectDays, criteria)

			if results.Total != tt.expectedCount {
				t.Errorf("Expected %d results for areas %v, got %d", tt.expectedCount, tt.areas, results.Total)
			}
		})
	}
}

func TestSearchByDateRange(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	tests := []struct {
		name          string
		dateFrom      string
		dateTo        string
		expectedCount int
	}{
		{"december only", "2023-12-01", "2023-12-31", 3},
		{"single day", "2023-12-01", "2023-12-01", 1},
		{"november only", "2023-11-01", "2023-11-30", 1},
		{"from december 2nd", "2023-12-02", "", 2},
		{"until december 1st", "", "2023-12-01", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := search.SearchCriteria{
				DateFrom: tt.dateFrom,
				DateTo:   tt.dateTo,
			}
			results := searchService.Search(perfectDays, criteria)

			if results.Total != tt.expectedCount {
				t.Errorf("Expected %d results for date range %s to %s, got %d",
					tt.expectedCount, tt.dateFrom, tt.dateTo, results.Total)
			}
		})
	}
}

func TestSearchSorting(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
		firstID   string
	}{
		{"date desc", "date", "desc", "id4"},
		{"date asc", "date", "asc", "id3"},
		{"title asc", "title", "asc", "id1"},
		{"title desc", "title", "desc", "id3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := search.SearchCriteria{
				SortBy:    tt.sortBy,
				SortOrder: tt.sortOrder,
			}
			results := searchService.Search(perfectDays, criteria)

			if len(results.PerfectDays) == 0 {
				t.Error("Expected at least one result")
				return
			}

			firstResult := results.PerfectDays[0]
			if firstResult.ID != tt.firstID {
				t.Errorf("Expected first result ID %s, got %s", tt.firstID, firstResult.ID)
			}
		})
	}
}

func TestSearchPagination(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	criteria := search.SearchCriteria{
		Limit:  2,
		Offset: 0,
	}
	results := searchService.Search(perfectDays, criteria)

	if results.Total != 4 {
		t.Errorf("Expected total 4, got %d", results.Total)
	}
	if len(results.PerfectDays) != 2 {
		t.Errorf("Expected 2 results in page, got %d", len(results.PerfectDays))
	}
	if results.Limit != 2 {
		t.Errorf("Expected limit 2, got %d", results.Limit)
	}
	if results.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", results.Offset)
	}

	criteriaPage2 := search.SearchCriteria{
		Limit:  2,
		Offset: 2,
	}
	resultsPage2 := searchService.Search(perfectDays, criteriaPage2)

	if len(resultsPage2.PerfectDays) != 2 {
		t.Errorf("Expected 2 results in page 2, got %d", len(resultsPage2.PerfectDays))
	}
	if resultsPage2.Offset != 2 {
		t.Errorf("Expected offset 2, got %d", resultsPage2.Offset)
	}
}

func TestSearchCombinedCriteria(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	criteria := search.SearchCriteria{
		Query:    "coffee",
		Username: "alice",
		Areas:    []string{"Shibuya"},
	}
	results := searchService.Search(perfectDays, criteria)

	if results.Total != 1 {
		t.Errorf("Expected 1 result for combined criteria, got %d", results.Total)
	}

	if len(results.PerfectDays) > 0 {
		result := results.PerfectDays[0]
		if result.Username != "alice" {
			t.Errorf("Expected username alice, got %s", result.Username)
		}
		if result.ID != "id1" {
			t.Errorf("Expected ID id1, got %s", result.ID)
		}
	}
}

func TestGetUniqueAreas(t *testing.T) {
	searchService := search.NewSearchService()
	perfectDays := createTestPerfectDays()

	areas := searchService.GetUniqueAreas(perfectDays)

	expectedAreas := []string{"Shibuya", "Shinjuku", "Ueno"}
	if len(areas) != len(expectedAreas) {
		t.Errorf("Expected %d areas, got %d", len(expectedAreas), len(areas))
	}

	for i, expected := range expectedAreas {
		if i >= len(areas) || areas[i] != expected {
			t.Errorf("Expected area %s at position %d, got %s", expected, i, areas[i])
		}
	}
}