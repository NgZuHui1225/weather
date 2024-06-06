package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
)

const apiKey = "L4V2B56VD6YY8KCJCJBB6DUSK"
const apiEndpoint = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

type WeatherResponse struct {
	Days []struct {
		Date          string  `json:"datetime"`
		Temperature   float64 `json:"temp"`
		Precipitation float64 `json:"precip"`
	} `json:"days"`
}

type QueryParams struct {
	Location  string `json:"location"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ParamData struct {
	data []QueryParams
}

func main() {
	r := chi.NewRouter()
	store := &ParamData{
		data: make([]QueryParams, 0),
	}

	client := resty.New()

	//POST METHOD
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST / endpoint hit")
		var param QueryParams
		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			log.Println("Error decoding JSON:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Received parameters: %+v\n", param)

		apiURL := fmt.Sprintf("%s%s/%s/%s", apiEndpoint, param.Location, param.StartDate, param.EndDate)
		params := url.Values{}
		params.Add("key", apiKey)
		params.Add("unitGroup", "metric")

		log.Printf("Requesting weather data from: %s\n", apiURL)

		resp, err := client.R().
			SetQueryParamsFromValues(params).
			SetHeader("Content-Type", "application/json").
			Get(apiURL)

		if err != nil {
			log.Println("Error making request: ", err)
			http.Error(w, "Error making request to weather API", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode() != 200 {
			log.Printf("Error: %v", resp)
			http.Error(w, "Error from weather API", resp.StatusCode())
			return
		}

		var weatherData WeatherResponse
		err = json.Unmarshal(resp.Body(), &weatherData)
		if err != nil {
			log.Println("Error decoding JSON response: ", err)
			http.Error(w, "Error decoding JSON response from weather API", http.StatusInternalServerError)
			return
		}

		log.Printf("Weather data: %+v\n", weatherData)

		store.data = append(store.data, param)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(weatherData); err != nil {
			log.Printf("Error encoding response: %v\n", err)
		}
	})

	//GET METHOD
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET / endpoint hit")
		data := store.data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	//localhost:3002
	log.Println("Starting server on :3002")
	if err := http.ListenAndServe(":3002", r); err != nil {
		log.Fatal(err)
	}

}
