# Centris-API: Kubernetes Implementation

A scalable web scraper for real estate data, optimized for Kubernetes deployment.

## Overview

This branch implements Kubernetes support for the Centris-API, enabling distributed and parallelized scraping jobs. It leverages Kubernetes' resource management and orchestration to efficiently distribute scraping workloads across multiple pods.

## Getting Started

### Prerequisites

-   Docker
-   Kubernetes cluster (or minikube/kind for local development)
-   kubectl CLI tool

### Kubernetes Deployment

Deploy the application to a Kubernetes cluster:

```bash
kubectl apply -f kind.yml
```

## Kubernetes Features

-   **Parallelized Scraping**: Configure the number of parallel pods for scraping
-   **Job Management**: Uses Kubernetes Jobs for reliable task completion
-   **Resource Optimization**: Efficiently distributes scraping workload
-   **Scalability**: Easily scale up or down based on workload requirements

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
PORT=8080
APP_ENV=local
DATABASE_URL=postgresql://username:password@localhost:5432/database_name
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgresql://username:password@localhost:5432/database_name
GOOSE_MIGRATION_DIR=./migrations
POD_INDEX=0
NUM_PODS=50
```

## Architecture

This implementation divides the scraping workload across multiple pods, with each pod responsible for a portion of the target data. The Kubernetes Job ensures all scraping tasks complete successfully and handles any pod failures automatically.
