postgres:
	docker run --name available -p 5000:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 -d postgres:17-alpine
createdb:
	docker exec -it available createdb --username=postgres --owner=postgres available
dropdb:
	docker exec -it available dropdb available
migratecreate:
	migrate create -ext sql -dir db/migration -seq "$(name)"
migrateup:
	migrate -path db/migration -database "postgres://postgres:12345@localhost:5000/available?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgres://postgres:12345@localhost:5000/available?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgres://postgres:12345@localhost:5000/available?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgres://postgres:12345@localhost:5000/available?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: createdb dropdb postgres migratecreate migrateup migrateup1 migratedown migratedown1 sqlc test server