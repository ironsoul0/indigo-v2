postgres:
	docker run --name indigo_postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it indigo_postgres createdb --username=root --owner=root indigo

dropdb:
	docker exec -it indigo_postgres dropdb indigo

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/indigo?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/indigo?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/indigo?sslmode=disable" -verbose down 1

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/indigo?sslmode=disable" -verbose up 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server migrateup1 migratedown1
