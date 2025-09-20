package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func PromptInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func PromptConfirm(prompt string) bool {
	response := PromptInput(prompt + " (y/N): ")
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func FormatTimeRange(startTime string, duration int) string {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return startTime
	}

	end := start.Add(time.Duration(duration) * time.Minute)
	return fmt.Sprintf("%s - %s", start.Format("15:04"), end.Format("15:04"))
}

func FormatDuration(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60

	if hours == 0 {
		return fmt.Sprintf("%dm", mins)
	}
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, mins)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func PadString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}