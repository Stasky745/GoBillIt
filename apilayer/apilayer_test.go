package apilayer

import (
	"fmt"
	"os"
	"testing"
)

func TestAPI(t *testing.T) {
	apiKey := os.Getenv("APILAYER_API_KEY")

	rate, err := GetRate(apiKey, "EUR", "USD")
	if err != nil {
		t.Fatalf("Error getting rate: %s", err)
	}

	fmt.Println(rate)
}
