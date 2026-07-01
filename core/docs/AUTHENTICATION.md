# Authentication & Security

Kobo uses a standard **Public Key / Private Secret** paradigm for authenticating Integrators against the Kobo API.

## API Key & Secret

When you provision a new Integrator (via the `/admin/integrators` internal endpoint), Kobo will generate and return a pair of credentials:
1. **API Key** (e.g. `kobo_live_pk_...` or `kobo_test_pk_...`)
2. **API Secret** (e.g. `kobo_live_sk_...` or `kobo_test_sk_...`)

### Environments (Live vs Sandbox)
The environment you are interacting with is strictly bound by the prefix of the keys you use.
- **Live Keys (`kobo_live_`)**: Will hit the production database and trigger real provisioning against Nomba's live systems.
- **Sandbox Keys (`kobo_test_`)**: Will hit the production database but trigger simulated, mock behavior or route to Nomba's sandbox endpoints.

**Important:** Both the API Key and the API Secret must match environments. Sending a `kobo_live_` key with a `kobo_test_` secret will result in an immediate `401 Unauthorized` without database lookups.

## Making Authenticated Requests

The Kobo API expects standard HTTP **Basic Authentication**.

- **Username**: Your API Key
- **Password**: Your API Secret

### Example Request (cURL)
```bash
curl -X POST https://api.kobo.triumphsystems.tech/v1/identities \
  -u kobo_test_pk_12345:kobo_test_sk_67890 \
  -H "Content-Type: application/json" \
  -d '{"external_reference": "user_123", "display_name": "John Doe"}'
```

## Security Best Practices
- **Never expose your API Secret**: The API Secret is hashed using SHA-256 in Kobo's database. We cannot retrieve it for you if lost. It must never be embedded in client-side code (browsers, mobile apps). 
- **Rotate compromised keys immediately**: Since the database only stores the SHA-256 hash, you can simply provision a new Integrator and migrate your integration if a secret is leaked.
- **Public Key visibility**: The API Key (`pk_`) is stored unhashed and acts as your identifier. It is generally safe to use as an identifier, but should still be treated with care.

## Nomba Webhook Authentication
For inbound events from Nomba, Kobo verifies the cryptographic integrity of the payload using **HMAC-SHA256**. 
Integrators do not need to worry about this. Kobo uses `NOMBA_WEBHOOK_SECRET` to validate the `X-Nomba-Signature` header, rejecting all spoofed or tampered events automatically.
