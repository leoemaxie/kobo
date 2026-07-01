---
sidebar_position: 2
title: Account Lifecycle
---

# The Account Lifecycle

Kobo ensures strict adherence to a finite state machine for all Identities. An identity can only transition through specific states, providing absolute safety against unauthorized operations.

## Lifecycle States

- **Pending**: The identity has been registered in Kobo, but the virtual account provisioning process is still communicating with Nomba.
- **Active**: The identity is fully provisioned. The virtual account is active and can receive and reconcile funds.
- **Limited**: The identity has hit a predefined compliance or volume limit (e.g., KYC ceiling). The account is active but payouts or specific operations might be restricted until upgraded.
- **Closing**: A termination request has been received. Inbound payments will now be rejected.
- **Closed**: The account is fully decommissioned. The underlying virtual account is deactivated at Nomba, and no further operations can occur.

## State Transitions
Transitions are handled entirely by the Kobo Core engine. For example, you cannot manually force an account from `Closed` to `Active` without hitting the dedicated `/reopen` endpoint which evaluates the legality of the transition.
