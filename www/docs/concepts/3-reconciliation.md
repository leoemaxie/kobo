---
sidebar_position: 3
title: Reconciliation Engine
---

# Reconciliation Engine

The Reconciliation Engine is the heart of Kobo. It guarantees that the internal ledger accurately reflects real-world bank transfers via the Monnify platform.

## How it works
1. **Webhook Reception**: Kobo receives a signed webhook from Monnify whenever a transfer hits a virtual account.
2. **Idempotency Check**: Kobo extracts the unique `transactionRef`. If this reference has been processed before, the webhook is immediately discarded. This protects against network retries and double-crediting.
3. **Signature Verification**: Kobo verifies the SHA-512 signature against your webhook secret.
4. **Attribution**: Kobo looks up the identity that owns the virtual account.
5. **Ledger Commit**: An inbound ledger entry is created atomically.

## Exceptions
If Kobo cannot attribute a transaction (e.g., the virtual account was deleted externally, or the webhook is malformed), it creates an **Exception** record. Exceptions must be resolved manually via the API or dashboard before the funds hit the user's ledger.
