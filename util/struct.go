package util

// Diagnostics
type Diagnostics struct {
	Countriesapi    int    `json:"countriesapi"`
	Meteoapi        int    `json:"meto_api"`
	Currencyapi     int    `json:"currency_api"`
	Notificationapi int    `json:"notification_db"`
	NumWebhooks     int    `json:"webhooks"`
	Version         string `json:"v1"`
	Uptime          int    `json:"starttime"`
}

// Registrations
type Registration struct {
	ID         string   `json:"id"`
	Country    string   `json:"country"`
	IsoCode    string   `json:"isoCode"`
	Features   Features `json:"features"`
	LastChange string   `json:"lastChange"`
}

// List of features included in registrations
type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

// Structs from the REST Countries API
type Country struct {
	CapitalCity          []string       `json:"capital"`
	LatitudeAndLongitude []float64      `json:"latlng"`
	Population           int            `json:"population"`
	Area                 float64        `json:"area"`
	Currencies           map[string]any `json:"currencies"`
}

// Structs from the Open Meteo API
type Weather struct {
	Hourly ForecastHourly `json:"hourly"`
}

type ForecastHourly struct {
	Temperature   []float64 `json:"temperature_2m"`
	Precipitation []float64 `json:"precipitation"`
}

// Structs from the Currencies API
type Currency struct {
	Rates map[string]float64 `json:"rates"`
}

// Structs for the Dashboard-endpoint
type DashboardResponse struct {
	Name          string            `json:"country"`
	Isocode       string            `json:"isoCode"`
	Features      DashboardFeatures `json:"features"`
	LastRetrieval string            `json:"lastRetrieval"`
}

type DashboardFeatures struct {
	Temperature      float64            `json:"temperature,omitempty"`
	Precipitation    float64            `json:"precipitation,omitempty"`
	Capital          string             `json:"capital,omitempty"`
	Coordinates      Coordinates        `json:"coordinates,omitempty"`
	Population       int                `json:"population,omitempty"`
	Area             float64            `json:"area,omitempty"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}
