#!/bin/bash

# Example: export DATABASE_URL="postgres://user:password@localhost:5432/football_db?sslmode=disable"
# Uncomment and edit the line below with your actual credentials
# export DATABASE_URL="postgres://user:password@localhost:5432/football_db?sslmode=disable"

go run cmd/web/main.go