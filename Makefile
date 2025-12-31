.PHONY: docker-up docker-down

docker-up:
	docker compose up -d

docker-down:
	docker compose down
