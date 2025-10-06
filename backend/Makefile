postgres:
	docker run --name available -p 5000:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 -d postgres:17-alpine
createdb:
	docker exec -it available createdb --username=postgres --owner=postgres available
dropdb:
	docker exec -it available dropdb available
migrateup:
	migrate -path db/migration -database "postgres://postgres:12345@localhost:5000/available?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://postgres:12345@localhost:5000/available?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: createdb dropdb postgres migrateup migratedown sqlc test