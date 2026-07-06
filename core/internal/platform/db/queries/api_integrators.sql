-- name: CreateApiIntegrator :one
INSERT INTO api_integrators (id, name, plan, status, production_access_granted)
VALUES ($1, $2, 'pay_as_you_go', 'active', false)
RETURNING *;

-- name: CreateApiCredential :one
INSERT INTO api_credentials (id, integrator_id, environment, key_id, secret_hash, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetApiIntegratorByKey :one
SELECT 
    i.id,
    i.name,
    c.secret_hash AS api_secret_hash,
    (c.environment = 'sandbox')::boolean AS is_sandbox
FROM api_integrators i
JOIN api_credentials c ON i.id = c.integrator_id
WHERE c.key_id = $1 AND c.revoked_at IS NULL;
