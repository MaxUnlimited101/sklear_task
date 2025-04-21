// internal/models/models.go
package models

type User struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Subscription struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	City      string `json:"city"`
	Condition string `json:"condition"`
}

type WeatherResponse struct {
	Weather []struct {
		Id          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
		Temp_min   float64 `json:"temp_min"`
		Temp_max   float64 `json:"temp_max"`
		Humidity   int     `json:"humidity"`
	} `json:"main"`
}

type SubscriptionDto struct {
	Email     string `json:"email"`
	City      string `json:"city"`
	Condition string `json:"condition"`
}
