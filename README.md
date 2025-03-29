# Project centris-api (MongoDB Migration)

This branch contains a version of the centris-api modified to work with MongoDB instead of the original SQL database.

## MongoDB Migration Details

The major changes include:

1. Refactored broker scraper to work with MongoDB
2. Refactored property scraper to work with MongoDB
3. Adjusted session management for Tor browser instances
4. Modified repository layer for MongoDB compatibility

The MongoDB schema is designed to maintain all the original functionality while leveraging document-oriented database benefits.

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
PORT=8080
APP_ENV=local
MONGODB_URI=mongodb://username:password@localhost:27017
MONGODB_DATABASE=centris
MONGODB_TIMEOUT=30
TOR_PROXY_HOST=localhost
TOR_PROXY_PORT=9050
TOR_CONTROL_HOST=localhost
TOR_CONTROL_PORT=9051
TOR_PASSWORD=your_tor_password
```

