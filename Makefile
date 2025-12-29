.PHONY: docker-up docker-down migrate-create migrate-up migrate-down run-user

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq -digit 6 $(NAME)

migrate-up:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) down

run-user:
	cd ./services/user && go run ./cmd/api
