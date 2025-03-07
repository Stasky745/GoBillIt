package utils

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func GetDate(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", s)
}

// formatDateString receives a format string and returns the formatted date
// according to the given rules for day, month, and year placeholders.
func FormatDate(date time.Time, format string) string {
	if format == "" {
		format = "YYYY-MM-DD"
	}

	// Define placeholders and values
	replacements := map[string]string{
		"YYYY": fmt.Sprintf("%d", date.Year()),
		"YY":   fmt.Sprintf("%02d", date.Year()%100),
		"MM":   fmt.Sprintf("%02d", int(date.Month())),
		"M":    fmt.Sprintf("%d", int(date.Month())),
		"DD":   fmt.Sprintf("%02d", date.Day()),
		"D":    fmt.Sprintf("%d", date.Day()),
		"m":    date.Month().String(),
		"d":    date.Weekday().String(),
	}

	// Sort placeholders by length (longest first) to prevent partial replacements
	keys := make([]string, 0, len(replacements))
	for key := range replacements {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	// Create a map for quick prefix lookups
	keyMap := make(map[string]struct{})
	for _, key := range keys {
		keyMap[key] = struct{}{}
	}

	// Optimized processing using a single pass
	var sb strings.Builder
	i := 0
	for i < len(format) {
		found := false
		for _, key := range keys {
			if i+len(key) <= len(format) && format[i:i+len(key)] == key {
				sb.WriteString(replacements[key]) // Append replacement
				i += len(key)                     // Skip past replaced token
				found = true
				break
			}
		}
		if !found {
			sb.WriteByte(format[i]) // Append character as-is
			i++
		}
	}

	return sb.String()
}
