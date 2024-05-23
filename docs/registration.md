# Registrations

The 'Registration' endpoint focuses on creating & managing dashboard configurations. These configurations
dictate how a dashboard is represented (weather, capital, population & other country-specific information) through the
'Dashboard' endpoint. 

Handles the following requests:
* **POST** - creates a new dashboard
* **GET** - retrieves a specific dashboard or all registered dashboards
* **PUT** - updates specified dashboard
* **DELETE** - deletes specified dashboard
* **PATCH** - updates specific fields for specified dashboard

## Endpoint
Endpoint for registration:
```
{{url}}/dashboard/v1/registration
```
* **{{url}}** is service's URL. 

## How It Works
The registrations system allows users to define and modify dashboard configurations that determine the data displayed on user dashboards.

## Register a new dashboard configuration
To register a new dashboard configuration, send a post request with the following settings:
### Request (POST) 

```
Method: POST
Path: /dashboard/v1/registrations/
Content type: application/json
```

#### Body (example):
```
{
   "country": "Norway",                                     // Indicates country name 
   "isoCode": "NO",                                         // Indicates two-letter ISO code for country 
   "features": {
                  "temperature": true,                      // Indicates whether temperature in degree Celsius is shown
                  "precipitation": true,                    // Indicates whether precipitation (rain, showers and snow) is shown
                  "capital": true,                          // Indicates whether the name of the capital is shown
                  "coordinates": true,                      // Indicates whether country coordinates are shown
                  "population": true,                       // Indicates whether population is shown
                  "area": true,                             // Indicates whether land area size is shown
                  "targetCurrencies": ["EUR", "USD", "SEK"] // Indicates which exchange rates (to target currencies) relative to the base currency of the registered country (in this case NOK for Norway) are shown
               }
}
```

### Response
Successful registration of a new dashboard configuration returns an ID for the configuration and the time 
it was last changed (will be curren time as it was just created).

#### Example response:
```
{
    "id": 123888388909032
    "lastChange": "2024-02-29 12:31"
}
```
* Content type: `application/json`
* Status code: 201 - status created on success, appropriate error message on fail. 

## View a specific registered dashboard configuration

Enables retrieval of a specific registered dashboard configuration by using its ID.

### Request (GET)
```
Method: GET
Path: /dashboard/v1/registrations/{id}
```

`{id}` is the ID associated with the specific configuration. \
Example request: ```/dashboard/v1/registrations/123888388909032```

### Response
Returns details of specified dashboard.

#### Example response:
```
{
   "id": 123888388909032,
   "country": "Norway",
   "isoCode": "NO",
   "features": {
                  "temperature": true,
                  "precipitation": true,
                  "capital": true,
                  "coordinates": true,
                  "population": true,
                  "area": true,
                  "targetCurrencies": ["EUR", "USD", "SEK"]
               },
    "lastChange": "20240229 14:07"
}
```
* Content type: `application/json`
* Status code: 200 - status ok on success, appropriate error message on fail.


## View **all registered dashboard configurations**

Lists all registered dashboard configurations. 

### Request (GET)
```
Method: GET
Path: /dashboard/v1/registrations/
```

### Response
Returns list of all registered dashboard configurations.

#### Example response:
```
[
   {
      "id": 123888388909032,
      "country": "Norway",
      "isoCode": "NO",
      "features": {
                     "temperature": true,
                     "precipitation": true,
                     "capital": true,
                     "coordinates": true,
                     "population": true,
                     "area": true,
                     "targetCurrencies": ["EUR", "USD", "SEK"]
                  }, 
      "lastChange": "20240229 14:07"
   },
   {
      "id": 18323883293923,
      "country": "Denmark",
      "isoCode": "DK",
      "features": {
                     "temperature": false,
                     "precipitation": true,
                     "capital": true,
                     "coordinates": true,
                     "population": false,
                     "area": true,
                     "targetCurrencies": ["NOK", "MYR", "JPY", "EUR"]
                  },
       "lastChange": "20240224 08:27"
   }
]
``` 
* Content type: `application/json`
* Status code: 200 - status ok on success, appropriate error message on fail.

## Replace a specific registered dashboard configuration
Enables the replacing of specific registered dashboard configuration by using its ID. Updates LastChange so when performing a "GET" on the same id afterward LastChange will represent the last modification of the dasbhoard configuration.

### Request (PUT)
```
Method: PUT
Path: /dashboard/v1/registrations/{id}
```
`id` is the ID associated with the specific configuration. \
Example request: ```/dashboard/v1/registrations/123888388909032```

#### Body (example): 
```
{
   "country": "Norway",
   "isoCode": "NO",
   "features": {
                  "temperature": false,              // this value is to be changed
                  "precipitation": true,
                  "capital": true,
                  "coordinates": true, 
                  "population": true,
                  "area": false,                     // this value is to be changed
                  "targetCurrencies": ["EUR", "SEK"] // this value is to be changed
               }
}
```
### Response
No body
* Status code: 200 - status ok on success, appropriate error message on fail.
* Body: empty

## Delete a specific registered dashboard configuration

Enabling the deletion of a specific registered dashboard configuration using its ID.

### Request (DELETE)
```
Method: DELETE
Path: /dashboard/v1/registrations/{id}
```
`id` is the ID associated with the specific configuration. \
Example request: ```/dashboard/v1/registrations/123888388909032```

### Response
No body
* Status code: 204 - status no content on success, appropriate error message on fail.
* Body: empty

## Update specific fields of a registered dashboard - additional feature

Enables selective updating of specified fields of a registered dashboard using its ID. Differentiates
from PUT as you can select *only* the specific fields you wish to update.

### Request (PATCH)
```
Method: PATCH
Path: /dashboard/v1/registrations/{id}
```

`id` is the ID associated with the specific configuration. \
Example request: ```/dashboard/v1/registrations/123888388909032```

#### Body (example):
**NOTE!** Unlike the other request bodies paths in PATCH has to start with capital letter as stored values are automatically capitalized in the databse.

So use: \
Country: "value" instead of country : "value" \
Features.Temperature :"value" instead of feature.temperature:"value"
```
{
   "Country":"Sweden",               //Field to be updated
   "IsoCode":"SE",                   //Field to be updated
   "Features.Temperature": true,     //Field to be updated
   "Features.Coordinates": true      //Field to be updated
}
```

### Response
* Status code: 200 - status ok on success, appropriate error message on fail.
* Body: empty
