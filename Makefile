sqlc:
	@sqlc generate

test:
	@go test -v -cover ./...


migrate_up:
	@migrate -path db/migrations -database "postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable" -verbose up
#mysql
#migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/simple_bank" -verbose up

migrate_up-docker:
	@docker run -v $(PWD)/db/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable -verbose up
#docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
 #    -path=/migrations/ -database postgres://localhost:5432/database up 2


.PHONY: sqlc test migrate_up migrate_up-docker
