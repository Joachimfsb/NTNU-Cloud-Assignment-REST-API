# Status : Monitoring service availability.

The status-endpoint checks the availability of different services used in our service, meaning the REST Countries API, Currency API, and Open Meteo API, but it also checks the availablity of the notification database. It also gives the exact number of webhooks that currently exists in the service.

It's important to note that if a stub is active, it will retrieve the status code for that stub-service instead of the real service, and if a stub isn't active it will check the status code of the real service.

## Endpoint
Endpoint for status:
```
{{url}}/dashboard/v1/status
```
* **{{url}}** is service's URL. 

Handles the following requests:
* **GET** - retrieves diagnostic data about the services used and our service.
## Example request:
### Request:
```
Method: GET
Path: /dashboard/v1/status
```
### Response:
* Content type: `application/json`
* Status code: 200 OK, appropriate error message on fail. 

Example body:
```
{
    "countriesapi": 200,
    "meto_api": 200,
    "currency_api": 200,
    "notification_db": 200,
    "webhooks": 2,
    "v1": "v1",
    "starttime": 10
}
```
