package apilayer

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	APILAYER_URL = "https://api.apilayer.com/exchangerates_data/latest?"
)

func GetRate(apiKey string, symbols string, base string) (float64, error) {
	requestUrl := APILAYER_URL + "symbols=" + symbols + "&base=" + base

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("apikey", apiKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result["rates"].(map[string]interface{})[symbols].(float64), nil
}
