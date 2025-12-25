-- name: CreateRestaurant :exec
INSERT INTO restaurants (id, name, address, phone, email, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
