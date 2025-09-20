package search

import (
	"perfect-day/pkg/models"
	"sort"
	"strings"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

type SearchCriteria struct {
	Query      string
	Areas      []string
	Username   string
	DateFrom   string
	DateTo     string
	SortBy     string
	SortOrder  string
	Limit      int
	Offset     int
}

type SearchResult struct {
	PerfectDays []*models.PerfectDay `json:"perfect_days"`
	Total       int                  `json:"total"`
	Offset      int                  `json:"offset"`
	Limit       int                  `json:"limit"`
}

func (ss *SearchService) Search(perfectDays []*models.PerfectDay, criteria SearchCriteria) *SearchResult {
	filtered := ss.filterPerfectDays(perfectDays, criteria)
	sorted := ss.sortPerfectDays(filtered, criteria.SortBy, criteria.SortOrder)

	total := len(sorted)
	start := criteria.Offset
	end := start + criteria.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	if criteria.Limit <= 0 {
		start = 0
		end = total
	}

	return &SearchResult{
		PerfectDays: sorted[start:end],
		Total:       total,
		Offset:      criteria.Offset,
		Limit:       criteria.Limit,
	}
}

func (ss *SearchService) filterPerfectDays(perfectDays []*models.PerfectDay, criteria SearchCriteria) []*models.PerfectDay {
	var filtered []*models.PerfectDay

	for _, pd := range perfectDays {
		if !ss.matchesCriteria(pd, criteria) {
			continue
		}
		filtered = append(filtered, pd)
	}

	return filtered
}

func (ss *SearchService) matchesCriteria(pd *models.PerfectDay, criteria SearchCriteria) bool {
	if criteria.Username != "" && pd.Username != criteria.Username {
		return false
	}

	if criteria.DateFrom != "" && pd.Date < criteria.DateFrom {
		return false
	}

	if criteria.DateTo != "" && pd.Date > criteria.DateTo {
		return false
	}

	if len(criteria.Areas) > 0 && !ss.matchesAreas(pd, criteria.Areas) {
		return false
	}

	if criteria.Query != "" && !ss.matchesQuery(pd, criteria.Query) {
		return false
	}

	return true
}

func (ss *SearchService) matchesAreas(pd *models.PerfectDay, areas []string) bool {
	for _, area := range areas {
		for _, pdArea := range pd.Areas {
			if strings.EqualFold(area, pdArea) {
				return true
			}
		}
	}
	return false
}

func (ss *SearchService) matchesQuery(pd *models.PerfectDay, query string) bool {
	searchableContent := pd.SearchableContent()
	queryLower := strings.ToLower(query)

	queryTerms := strings.Fields(queryLower)
	for _, term := range queryTerms {
		if !strings.Contains(searchableContent, term) {
			return false
		}
	}

	return true
}

func (ss *SearchService) sortPerfectDays(perfectDays []*models.PerfectDay, sortBy, sortOrder string) []*models.PerfectDay {
	sorted := make([]*models.PerfectDay, len(perfectDays))
	copy(sorted, perfectDays)

	switch sortBy {
	case "date":
		sort.Slice(sorted, func(i, j int) bool {
			if sortOrder == "desc" {
				return sorted[i].Date > sorted[j].Date
			}
			return sorted[i].Date < sorted[j].Date
		})
	case "created_at":
		sort.Slice(sorted, func(i, j int) bool {
			if sortOrder == "desc" {
				return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
			}
			return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
		})
	case "title":
		sort.Slice(sorted, func(i, j int) bool {
			if sortOrder == "desc" {
				return sorted[i].Title > sorted[j].Title
			}
			return sorted[i].Title < sorted[j].Title
		})
	default:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		})
	}

	return sorted
}

func (ss *SearchService) GetUniqueAreas(perfectDays []*models.PerfectDay) []string {
	areaSet := make(map[string]bool)

	for _, pd := range perfectDays {
		for _, area := range pd.Areas {
			areaSet[area] = true
		}
	}

	areas := make([]string, 0, len(areaSet))
	for area := range areaSet {
		areas = append(areas, area)
	}

	sort.Strings(areas)
	return areas
}