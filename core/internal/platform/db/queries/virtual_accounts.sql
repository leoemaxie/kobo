-- name: CreateVirtualAccount :one
INSERT INTO virtual_accounts (id, identity_id, nomba_account_ref, account_number, bank_name, account_name, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetVirtualAccountByID :one
SELECT * FROM virtual_accounts
WHERE id = $1;

-- name: GetActiveVirtualAccountByIdentityID :one
SELECT * FROM virtual_accounts
WHERE identity_id = $1 AND is_active = true;

-- name: GetVirtualAccountByAccountNumber :one
SELECT * FROM virtual_accounts
WHERE account_number = $1;

-- name: UpdateVirtualAccountProvisioning :one
-- Called after a successful provisioning to set the account number and bank info
UPDATE virtual_accounts
SET account_number = $2,
    bank_name = $3,
    account_name = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateVirtualAccount :exec
-- Sets is_active to false when a new account is provisioned for the identity
UPDATE virtual_accounts
SET is_active = false,
    updated_at = now()
WHERE identity_id = $1 AND is_active = true;

-- name: ListActiveVirtualAccounts :many
SELECT * FROM virtual_accounts
WHERE is_active = true
ORDER BY id ASC;
