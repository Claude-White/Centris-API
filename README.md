# Centris-API

A Go-based API service for scraping and serving real estate data from Centris.ca, focused on Quebec real estate market listings.

## Project Overview

Centris-API is a powerful tool that extracts, processes, and serves real estate data from Centris.ca through a RESTful API. It manages two main types of data:

1. **Properties** - Real estate listings with details like price, location, features, and photos
2. **Brokers** - Real estate agents with their contact information, served areas, and associated properties

The system includes two scrapers:

-   **Property Scraper** - Collects property information including expenses, features, and photos
-   **Broker Scraper** - Gathers broker details, phone numbers, and external links (social media, websites)

## Features

-   **RESTful API** with Swagger documentation
-   **Property data** including features, expenses, and photos
-   **Broker information** with contact details and external links
-   **Geolocation-based search** capabilities
-   **Pagination support** for all endpoints
-   **PostgreSQL database** for data storage
-   **Docker support** for easy deployment
-   **Graceful shutdown** handling

## Technology Stack

-   **Go** (1.23+) - Core programming language
-   **PostgreSQL** - Database for storing scraped data
-   **SQLC** - Type-safe SQL query generation
-   **Swagger** - API documentation
-   **Docker/Docker Compose** - Containerization and orchestration
-   **Goose** - Database migration management

## API Endpoints

### Properties

-   `GET /properties/{mls}` - Get property by MLS number
-   `POST /properties` - Get all properties (paginated)
-   `POST /properties/coordinates` - Get property by coordinates
-   `POST /properties/broker` - Get all properties by broker
-   `POST /properties/agency` - Get all properties by agency
-   `POST /properties/city` - Get all properties by city
-   `POST /properties/radius` - Get all properties within a radius

### Brokers

-   `GET /brokers/{brokerId}` - Get broker by ID
-   `POST /brokers` - Get all brokers (paginated)

## Getting Started

### Prerequisites

-   Go 1.23 or higher
-   PostgreSQL
-   Docker and Docker Compose (optional)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
PORT=8080
APP_ENV=local
DATABASE_URL=postgresql://username:password@localhost:5432/database_name
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgresql://username:password@localhost:5432/database_name
GOOSE_MIGRATION_DIR=./migrations
```

### Installation

1. Clone the repository

    ```bash
    git clone https://github.com/your-username/Centris-API.git
    cd Centris-API
    ```

2. Install dependencies

    ```bash
    go mod download
    ```

3. Setup the database

4. Build and run the application

    ```bash
    go run cmd/api/main.go [property-scraper|broker-scraper]
    ```

5. Access the Swagger documentation at http://localhost:8080/swagger/

## Command Line Options

The application can be run with different commands:

-   Default (no args) - Start the HTTP server
-   `broker-scraper` - Run the broker data scraper
-   `property-scraper` - Run the property data scraper

Example:

```bash
./main property-scraper
```

## Project Structure

-   `/cmd/api` - Application entrypoint and main server code
-   `/docs` - Swagger documentation
-   `/internal`
    -   `/repository` - Database models and queries (SQLC-generated)
    -   `/server` - HTTP server, routes, and scraper implementations
-   `/migrations` - Database migration files
-   `/sqlc` - SQLC configuration and query definitions
    -   `/queries` - SQL query templates

## Branches

-   `main` - Stable version of the application
-   `MongoDB-Migration` - Migration to MongoDB database
-   `kubernetes-implementation` - Kubernetes deployment configuration

## Acknowledgments

-   [Centris.ca](https://www.centris.ca) - Source of real estate data
-   [SQLC](https://sqlc.dev/) - SQL query generation tool
-   [Swagger](https://swagger.io/) - API documentation framework
-   [Goose](https://github.com/pressly/goose) - Database migration tool
