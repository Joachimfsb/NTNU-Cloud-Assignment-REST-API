# Firebase
Firebase is Google's database program which is a NoSQL type of database. This is used to store the data of users so they wont lose different things they have saved (in this instance dashboard configurations)

# Content
- [Overview](#overview)
- [Usage](#usage)

## Overview
The database is split into one collection ('dashboards'), and multiple documents with unique IDs. See what the documents contains [here](#document-contents)

## Usage
To use this service, the database needs a collection called 'dashboards'. This collection is where the dashboard configurations of the users will be stored. 
The 'dashboards' collection is where the configurations is stored, and these configurations will automatically get unique IDs that firebase creates.

### Collections
As mentioned earlier, there is one collection named 'dashboards' that contains all the documents. Other collections will not be compatible with this service.

### Documents
These are what actually saves/stores the different configurations the user wishes to store. There is no limit to how many configurations/documents a user can store.

Every documentation has a unique ID consisting of capital letters, lowercase letter and numbers. This makes it harder for unauthorized entites to guess the IDs of the different configurations stored.

### Document contents:
- country (string)
- isoCode (string)
- features (map)
    - temperature (bool)
    - precipitation (bool)
    - capital (bool)
    - coordinates (bool)
    - population (bool)
    - area (bool)
    - targetCurrencies (array)
- lastChanged (string)

### firebaseKey.json
To connect to our database, you need to have access to the key for this database. This file needs to be in the root directory and never to be committed to Git!
