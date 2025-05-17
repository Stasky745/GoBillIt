package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Stasky745/GoBillIt/internal/utils"
	"github.com/Stasky745/go-libs/log"
)

// SetDefaultEnv is a generic function that returns a value of any type based on an environment variable
func setDefaultEnv[T any](env string, def T) T {
	if v, ok := os.LookupEnv(env); ok {
		switch any(def).(type) {
		case string:
			if v == "" {
				return def
			}
			return any(v).(T)
		case bool:
			// Handle the case where the environment variable is a bool
			if strings.ToLower(v) == "false" || v == "0" {
				return any(false).(T)
			}
			if strings.ToLower(v) == "true" || v == "1" {
				return any(true).(T)
			}
		case int:
			if intValue, err := strconv.Atoi(v); err == nil {
				return any(intValue).(T)
			}
		case int64:
			if intValue, err := strconv.ParseInt(v, 10, 64); err == nil {
				return any(intValue).(T)
			}
		case float64:
			if floatValue, err := strconv.ParseFloat(v, 64); err == nil {
				return any(floatValue).(T)
			}
		}
	}

	// Return the default value if the environment variable is not found or cannot be converted
	return def
}

// replaceTemplateValues searches for "{{ key }}" in text and replaces with k.String("key")
func template(text string, m map[string]string) string {
	// Regex pattern to find {{ key[aux] }} placeholders
	re := regexp.MustCompile(REGEX_TEMPLATING)

	// Replace function: Fetch the value from Koanf
	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract key and aux value (if any)
		parts := re.FindStringSubmatch(match)
		if parts == nil {
			return text
		}
		key := strings.TrimSpace(parts[1])
		aux := strings.TrimSpace(parts[2])

		switch key {
		case "date":
			if aux == "" {
				aux = "YYYYMM"
			}
			return utils.FormatDate(time.Now(), aux)
		case "inv.date", "inv.duedate":
			date, err := utils.GetDate(k.String(key))
			if log.CheckErr(err, false, "can't parse date, will use today", "key", key, "date", k.String(key)) {
				date = time.Now()
			}
			return utils.FormatDate(date, aux)
		case "conversion", "inv.conversion", "inv.conversion.value":
			return strconv.FormatFloat(k.Float64("inv.conversion.value"), 'f', -1, 64)
		case "inv.conversion.min":
			return strconv.FormatFloat(k.Float64("inv.conversion.min"), 'f', -1, 64)
		default:
			if val, ok := m[key]; ok {
				return val
			}
			return template(k.String(key), m)
		}
	})

	return result
}

func templateExtractParams(text, s string) string {
	// Regex pattern to find {{ key[aux] }} placeholders
	re := regexp.MustCompile(REGEX_TEMPLATING)

	// ReplaceAllStringFunc to process the text and resolve templates recursively
	res := re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract key and aux value (if any)
		parts := re.FindStringSubmatch(match)
		if parts == nil {
			return text
		}
		templateKey := strings.TrimSpace(parts[1])
		aux := strings.TrimSpace(parts[2])

		// If the key matches the target key, return the aux value (if exists)
		if templateKey == s {
			return "{{{" + aux + "}}}"
		}

		// Fetch the value from Koanf (e.g., k.String or similar)
		// Check if the key exists in Koanf
		val := k.String(templateKey)
		if val == "" {
			// If key does not exist, return the original match (to avoid infinite loops)
			return match
		}

		// Otherwise, recursively resolve the value of the key (template resolution)
		return templateExtractParams(val, s)
	})

	return res
}

func templateGetKeyParams(text, s string) (bool, string) {
	newText := templateExtractParams(text, s)
	// Regex pattern to find {{ key[aux] }} placeholders
	re := regexp.MustCompile(`\{\{\{\s*(\w+)\s*\}\}\}`)
	res := re.FindStringSubmatch(newText)
	if res == nil {
		return false, ""
	} else if len(res) < 2 {
		return true, ""
	}

	return true, strings.TrimSpace(res[1])
}
