.PHONY: migrate-create migrate-up migrate-down server

migrate-create:
	migrate create -ext sql -dir migrations -seq -digits 6 $(NAME)

migrate-up:
	migrate -path migrations -database "postgres://root:secret@localhost:5432/airbnb?sslmode=disable&x-migrations-table=user_schema_migrations" up

migrate-down:
	migrate -path migrations -database "postgres://root:secret@localhost:5432/airbnb?sslmode=disable&x-migrations-table=user_schema_migrations" down

server:
	go run ./cmd/api
