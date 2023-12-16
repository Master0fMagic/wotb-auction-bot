package dto

import "time"

// Price represents the price details in the JSON response
type Price struct {
	Currency struct {
		Name     string   `json:"name"`
		Count    int      `json:"count"`
		Title    string   `json:"title"`
		ImageURL string   `json:"image_url"`
		Sizes    struct{} `json:"sizes"`
		Type     string   `json:"type"`
	} `json:"currency"`
	Value int `json:"value"`
}

// Entity represents the entity details in the JSON response
type Entity struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Nation           string `json:"nation"`
	Subnation        string `json:"subnation"`
	UseSubnationFlag bool   `json:"use_subnation_flag"`
	TypeSlug         string `json:"type_slug"`
	Level            int    `json:"level"`
	RomanLevel       string `json:"roman_level"`
	UserString       string `json:"user_string"`
	ImageURL         string `json:"image_url"`
	PreviewImageURL  string `json:"preview_image_url"`
	IsPremium        bool   `json:"is_premium"`
	IsCollectible    bool   `json:"is_collectible"`
}

// Result represents each item in the "results" array in the JSON response
type Result struct {
	ID                 int    `json:"id"`
	Type               string `json:"type"`
	Countable          bool   `json:"countable"`
	Entity             Entity `json:"entity"`
	InitialCount       int    `json:"initial_count"`
	CurrentCount       int    `json:"current_count"`
	Saleable           bool   `json:"saleable"`
	AvailableFrom      string `json:"available_from"`
	AvailableBefore    string `json:"available_before"`
	Price              Price  `json:"price"`
	NextPrice          *Price `json:"next_price"`
	Available          bool   `json:"available"`
	Display            bool   `json:"display"`
	NextPriceDatetime  string `json:"next_price_datetime"`
	NextPriceTimestamp int    `json:"next_price_timestamp"`
}

// Response represents the entire JSON response
type Response struct {
	Count   int      `json:"count"`
	HasNext bool     `json:"has_next"`
	Results []Result `json:"results"`
}

// Data represents the data structure from the API
type Data struct {
	Entities []Entity `json:"entities"`
}

// Vehicle represents a vehicle with additional information
type Vehicle struct {
	ID                 int
	Type               string
	Countable          bool
	Entity             Entity
	InitialCount       int
	CurrentCount       int
	Saleable           bool
	AvailableFrom      time.Time
	AvailableBefore    time.Time
	Price              Price
	NextPrice          *Price
	Available          bool
	Display            bool
	NextPriceDatetime  time.Time
	NextPriceTimestamp int
}
