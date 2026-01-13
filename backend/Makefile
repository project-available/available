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
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/project-available/available.git/db/sqlc Store
sqlc:
	sqlc generate
test:
	go test -v -coverprofile=coverage.out ./...
report:
	go tool cover -html=coverage.out -o coverage.html
server:
	go run main.go
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	swag init
swagger-clean:
	rm -rf docs/
.PHONY: createdb dropdb postgres migratecreate migrateup migrateup1 migratedown migratedown1 sqlc test server mock swagger-install swagger swagger-clean