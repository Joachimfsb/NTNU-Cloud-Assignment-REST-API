# Dashboards : Retrieve populated dashboard.

This endpoint can be used to retrieve and display information about a specific registration.

The dashboard-endpoint uses three services:

    - REST Countries API
        - Endpoint: http://129.241.150.113:8080/v3.1
        - Documentation: http://129.241.150.113:8080/

    - Open-Meteo API
        - Documentation: https://open-meteo.com/en/features#available-apis

    - Currency API
        - Endpoint: http://129.241.150.113:9090/currency/
        - Documentation: http://129.241.150.113:9090/

One or all of these services could be replaced by a stub-service which can be activated through the 
"config.yaml"-file, by changing the wanted stub-services value to `true`. These stubs returns mocked
JSON-data from the original/real service, meaning it returns static data for Norway gotten and stored at a
specific time. Instead of dynamically retrieved data from a real service.

## Endpoint
Endpoint for dashboards:
```
{{url}}/dashboard/v1/dashboards
```
* **{{url}}** is service's URL. 

Handles the following requests:
* **GET** - retrieves a specific dashboard with it's wanted data
## Example request:
### Request:
```
Method: GET
Path: /dashboard/v1/dashboards/<id>
```
* **id** is a specific ID of a registration which you can find using the registration-endpoint.. 

### Response:
* Content type: `application/json`
* Status code: 200 OK, appropriate error message on fail. 

Example body:
```
{
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "temperature": -1.9035714285714294,
        "precipitation": 0.027380952380952384,
        "capital": "Oslo",
        "coordinates": {
            "latitude": 62,
            "longitude": 10
        },
        "population": 5379475,
        "area": 323802,
        "targetCurrencies": {
            "EUR": 0.085272,
            "SEK": 0.995781,
            "USD": 0.090918
        }
    },
    "lastRetrieval": "2024-04-18 17:35"
}
```
