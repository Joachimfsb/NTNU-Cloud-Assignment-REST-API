package util

const (
	// Endpoints
	REGISTRATION_PATH = "/dashboard/v1/registrations/"
	DASHBOARD_PATH    = "/dashboard/v1/dashboards/"
	NOTIFICATION_PATH = "/dashboard/v1/notifications/"
	STATUS_PATH       = "/dashboard/v1/status/"

	// URLs
	COUNTRY_URL   = "http://129.241.150.113:8080/v3.1"
	OPEN_METO_URL = "https://api.open-meteo.com/v1/forecast"
	CURRENCY_URL  = "http://129.241.150.113:9090/currency/"
	LOCALHOST     = "http://localhost:"

	// Collections
	DASHBOARDS               = "dashboards"
	COLLECTION_NOTIFICATIONS = "notifications"

	// STUB Ports
	DATABASE_PORT       = "1881"
	CURRENCIES_PORT     = "13272"
	OPENWEATHER_PORT    = "17623"
	REST_COUNTRIES_PORT = "25531"
)

var (
	DatabaseStubPort string = DATABASE_PORT

	WeatherStubPort    string = OPENWEATHER_PORT
	CurrenciesStubPort string = CURRENCIES_PORT
	CountryStubPort    string = REST_COUNTRIES_PORT

	// Stubs
	STUB_DATABASE_REGISTRATIONS = "stubs/res/registrations.json"
	STUB_DATABASE_NOTIFICATIONS = "stubs/res/notifications.json"

	// Mocked response from the real services
	STUB_WEATHER_REPONSE     = "stubs/res/weather.json"
	STUB_CURRENCIES_RESPONSE = "stubs/res/currency.json"
	STUB_COUNTRY_RESPONSE    = "stubs/res/country.json"
)
