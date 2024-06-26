# Dashboard Configurations
Create your very own dashboards and configure them just the way you like it! *Dashboard Configurations* is a Machine-to-Machine API that lets you configure dashboards with currency exchance rates, weather data and more per country.

Note: The documentation you can read from this README.md file and other Markdown files in the [docs](./docs) directory are meant for internal use only. Public documentation for this project does not and will not exist. The simple reason is this is an assignment.

## Navigations
- [Installation](#installation)
- [Usage](#usage)
- [Read More](#read-more)

## Installation
To get started, you need [Go](https://go.dev/) version 1.21. Head over to the [downloads](https://go.dev/dl/) page to find the binary for your system.

We also recommend to install [Git](https://git-scm.com/). On *NIX systems, this is usually installed by default, or you can use your system's package manager to install it. For NT, head over to [Git's downloads page](https://git-scm.com/download/win) to find the binary for your system.

The project uses Google's Cloud Firestore as DBMS. Head over to your Firestore instance and download the public key.

Clone the project to your system:
```
git clone https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2024-workspace/joafsb/assignment2.git
cd assignment2
```

Move your public key to the project. Now you need to edit the configuration file and insert the key's location. The configuration is [config.yaml](./config.yaml). Under secrets/firebaseKey, insert the file path to the key. It supports both absolute and relative paths.

Run the project
```
go run main.go
```

The project has Docker support. Build and run the container as follows:
```
docker-compose build
docker run -d -p 0.0.0.0:8080:8080 assignment2-demoapp
```

You can also expose the port only to localhost as follows:
```
docker run -d -p 127.0.0.1:8080:8080 assignment2-demoapp
```

## Usage
The project consists of the following four modules:
- [Dashboards](./docs/dashboards.md)
- [Registrations](./docs/registration.md)
- [Notifications](./docs/notifications.md)
- [Status](./docs/status.md)

Read the API documentations by navigating to the above links.

## Read More
Below are links to more documentations about the implementation of this project's components. Read more about them by navigating to the below links:
- [Firebase](./docs/firebase.md)
