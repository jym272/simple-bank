sqlc:
	@sqlc generate

test:
	@go test -v -cover ./...


migrate_up:
	@migrate -path db/migrations -database "postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable" -verbose up
#mysql
#migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/simple_bank" -verbose up


.PHONY: sqlc test migrate_up
