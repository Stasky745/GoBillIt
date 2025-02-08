package utils

import (
	"fmt"
	"math"
	"strings"
)

func GetConvertedCost(price float64, conversion float64) float64 {
	return math.Round(price*conversion*100) / 100
}

func FormatFloatToAmount(amount float64) string {
	// Format the float to two decimal places
	formatted := fmt.Sprintf("%.2f", amount)

	// Split the whole part and the decimal part
	parts := strings.Split(formatted, ".")
	wholePart := parts[0]
	decimalPart := ""

	// Handle decimal part if it exists
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}

	// Add commas to the whole part
	commaSeparated := addCommas(wholePart)

	// Combine the whole part and the decimal part
	return commaSeparated + decimalPart
}

func addCommas(s string) string {
	// Reverse the string
	reversed := reverseString(s)
	var sb strings.Builder

	for i, ch := range reversed {
		// Add a comma after every 3 digits, except for the last group
		if i > 0 && i%3 == 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune(ch)
	}

	// Reverse the string again to get the correct order
	return reverseString(sb.String())
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
