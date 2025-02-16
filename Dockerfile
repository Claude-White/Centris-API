FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod

RUN apk add --no-cache tor

WORKDIR /app
COPY --from=build /app/main /app/main

ENV PORT=8080
ENV APP_ENV=local
ENV DATABASE_URL=postgresql://postgres.wozrsuxaxijneojygblo:contests-whisk-borrow-PLUNDER3@aws-0-ca-central-1.pooler.supabase.com:5432/postgres?sslmode=require&options=-c%20statement_timeout%3D120000
ENV POD_INDEX=0
ENV NUM_PODS=100

EXPOSE ${PORT}
CMD tor & ./main property-scraper