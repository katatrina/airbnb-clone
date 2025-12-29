.PHONY: docker-up docker-down run-user

docker-up:
	docker compose up -d

docker-down:
	docker compose down

run-user:
	cd ./services/user && go run ./cmd/api
