# Project centris-api (MongoDB Migration)

This branch contains a version of the centris-api modified to work with MongoDB instead of the original SQL database.

## MongoDB Migration Details

The major changes include:

1. Refactored broker scraper to work with MongoDB
2. Refactored property scraper to work with MongoDB
3. Adjusted session management for Tor browser instances
4. Modified repository layer for MongoDB compatibility

The MongoDB schema is designed to maintain all the original functionality while leveraging document-oriented database benefits.

