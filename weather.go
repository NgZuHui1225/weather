package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

const apiKey = "FZ4UNKA5K4BXQUBLYWT3NLT2A"
const apiEndpoint = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s/%s?unitGroup=metric&key=%s"

type WeatherResponse struct {
	Days []struct {
		Date          string  `json:"datetime"`
		Temperature   float64 `json:"temp"`
		Precipitation float64 `json:"precip"`
		// Add more fields as needed
	} `json:"days"`
}

func main() {
	location := "Kuala Lumpur, KL"
	startDate := "2024-06-01"
	endDate := "2024-06-03"

	client := resty.New()

	url := fmt.Sprintf(apiEndpoint, location, startDate, endDate, apiKey)

	resp, err := client.R().
		//SetQueryParam("key", apiKey).
		//SetQueryParam("location", location).
		//SetQueryParam("startDate", startDate).
		//SetQueryParam("endDate", endDate).
		SetHeader("Content-Type", "application/json").
		//Get(apiEndpoint)
		Get(url)

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Fatalf("Error: %v", resp)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(resp.Body(), &weatherData)
	if err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	// Print weather data
	fmt.Printf("Weather for %s from %s to %s:\n", location, startDate, endDate)
	for _, day := range weatherData.Days {
		fmt.Printf("Date: %s, Temperature: %.2fC, Precipitation: %.2fmm\n", day.Date, day.Temperature, day.Precipitation)
	}
}
