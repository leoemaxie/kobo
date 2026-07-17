---
sidebar_position: 1
title: Core Primitives
---

# Core Primitives

Kobo is built around three fundamental concepts: **Identities**, **Virtual Accounts**, and **Ledgers**.

## Identities
An `Identity` is the root entity in Kobo. It represents an end-user, a business, or an entity that requires financial tracking. 
Identities have a strictly managed lifecycle (Pending, Active, Limited, Closing, Closed) that controls what financial operations are permitted.

## Virtual Accounts
Every Identity is backed by a `Virtual Account`. These are real, dedicated NUBAN accounts provisioned via Monnify. 
- **1-to-1 Mapping**: An identity has exactly one active virtual account at a time.
- **Auto-Provisioning**: When an identity is created, Kobo automatically orchestrates the provisioning of the virtual account.

## Ledgers
The `Ledger` is the immutable source of truth for all balances. Kobo uses an append-only transaction system.
- **Inbound Entries**: Represent money credited to the virtual account via a reconciled bank transfer.
- **Outbound Entries**: Represent money moving out of the system (e.g., payouts or sweeps).

Your application never calculates balances manually. Instead, it queries the Ledger to get the computed balance for a given Identity at any point in time.
