.PHONY: migrate-create migrate-up migrate-down server

migrate-create:
	migrate create -ext sql -dir migrations -seq -digits 6 $(NAME)

migrate-up:
	migrate -path migrations -database "postgres://root:secret@localhost:5432/airbnb?sslmode=disable&x-migrations-table=booking_schema_migrations" -verbose up

migrate-down:
	migrate -path migrations -database "postgres://root:secret@localhost:5432/airbnb?sslmode=disable&x-migrations-table=booking_schema_migrations" -verbose down

server:
	go run ./cmd/api
