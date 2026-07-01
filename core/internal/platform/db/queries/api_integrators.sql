-- name: CreateApiIntegrator :one
INSERT INTO api_integrators (id, name, api_key, api_secret_hash, is_sandbox)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetApiIntegratorByKey :one
SELECT * FROM api_integrators
WHERE api_key = $1;
