#!/usr/bin/env bash


migrate create -ext sql -dir db/migrations -seq init_schema

migrate -path db/migrations -database "postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable" -verbose up

migrate -path db/migrations -database "postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable" -verbose down
# psql
# see all db
