# Centris-API: Kubernetes Implementation

A scalable web scraper for real estate data, optimized for Kubernetes deployment.

## Overview

This branch implements Kubernetes support for the Centris-API, enabling distributed and parallelized scraping jobs. It leverages Kubernetes' resource management and orchestration to efficiently distribute scraping workloads across multiple pods.

## Getting Started

### Prerequisites

- Docker
- Kubernetes cluster (or minikube/kind for local development)
- kubectl CLI tool

### Kubernetes Deployment

Deploy the application to a Kubernetes cluster:

```bash
kubectl apply -f kind.yml
```

## Kubernetes Features

- **Parallelized Scraping**: Configure the number of parallel pods for scraping
- **Job Management**: Uses Kubernetes Jobs for reliable task completion
- **Resource Optimization**: Efficiently distributes scraping workload
- **Scalability**: Easily scale up or down based on workload requirements

## Environment Variables

The application requires the following environment variables:

- `APP_ENV`: Application environment (development/production)
- `PORT`: Application port
- `DATABASE_URL`: PostgreSQL database connection string
- `GOOSE_DRIVER`: Database driver for Goose migrations (usually "postgres")
- `GOOSE_DBSTRING`: Database connection string for Goose migrations
- `GOOSE_MIGRATION_DIR`: Directory containing migration files
- `NUM_PODS`: Total number of pods for distributed scraping
- `POD_INDEX`: Pod index for workload distribution (automatically set by Kubernetes)

## Architecture

This implementation divides the scraping workload across multiple pods, with each pod responsible for a portion of the target data. The Kubernetes Job ensures all scraping tasks complete successfully and handles any pod failures automatically.
