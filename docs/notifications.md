# Notifications
Our notifications system is a webhook service. It allows clients to subscribe to events. The following events are available:

- **REGISTER** - notification is sent if a new configuration is registered
- **CHANGE** - notification is sent if configuration is modified
- **DELETE** - notification is sent if configuration is deleted
- **INVOKE** - notification is sent if dashboard is retrieved (i.e., populated with values)

## Endpoint
The endpoint to the notifications is as follows:
```
{{url}}/dashboard/v1/notifications
```

Details:
- `{{url}}` is the service's URL.

## How It Works
Whenever an event happens, notifications are sent to whomever subscribed to them. For example, when a new dashboard configuration is registered, a notification is sent to all subscribers to signed up to the configuration's country code.

## Subscribe to an Event
Registering to an event is simple. Send a POST request with the following request body to get started:

```json
{
    "Url": "https://your-url.no",
    "Country": "COUNTRY-CODE",
    "Event": "EVENT-TYPE"
}
```

Details:
- `Url` is the URL we will POST to whenever an event occurs.
- `Country` is the country code which the event will invoke on. For example, if the notification subscribes to "NO", then an event that configures that specific country code will notify the subscriber. Please note that this is a country code. The characters must be all uppercase.
- `Event` is the type of event to notify about. If a notification subscribes to for example **INVOKE**, then an invoke event will send a notification to the registered URL.

### Example
```json
{
    "url": "https://api.my-notifications.test",
    "country": "NO",
    "event": "REGISTER"
}
```

Output (random ID):
```json
{
    "id": "17388271632813271"
}
```

### Output
**Success**: If a new notification was successfully registered, then its new ID is returned in the body and the status code is *201 Created*
**Error**: If a field is missing, then the server returns *422 Unprocessable* and a message about which field(s) is missing.

## Delete a Notification
Notifications can be deleted if its id is passed to the request as URl parameter.

### Example
```http
{{url}}/dashboard/v1/notifications/{{id}}
```

Details:
- `{{url}}` is the service's url
- `{{id}}` is the notification's id

### Output
Both success and error will return 204 No Content. This request is idempotent and does not guarantee that a notification is deleted.

## View All Notifications
To view all registered notifications, send a GET request to the [endpoint](#endpoint).

### Example
```
{{url}}/dashboard/v1/notifications
```

Details:
- {{url}} is the service's URL

Output:
```json
[
    {
        "url": "https://api.my-notifications.test",
        "country": "NO",
        "event": "REGISTER"
    },
    {
        "url": "https://api.some-notification.test",
        "country": "GB",
        "event": "INVOCATION"
    }
]
```

### Output
There are two possible outcomes:
1. Notifications exist. All notifications will be output in an array of JSON objects. Status code is 200 OK
2. No notifications exist. An empty JSON array is returned and the status code is 200 OK.

## View One Notifications
To view one registered notifications, send a GET request to the [endpoint](#endpoint) with the id.

### Example
```
{{url}}/dashboard/v1/notifications/{{id}}
```

Details:
- `{{url}}` is the service's URL.
- `{{id}}` is the notification's id.

Output:
```json
{
    "url": "https://api.my-notifications.test",
    "country": "NO",
    "event": "REGISTER"
}
```

### Output
There are two possible outcomes:
1. Notifications exist. All notifications will be output in an array of JSON objects. Status code is 200 OK.
2. The notification does not exist by the given ID. Status code is 404 Page Not Found.
