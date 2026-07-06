# Background Worker and Scheduled Jobs

The Kobo Core platform utilizes a long-running Go background worker process to handle scheduled tasks ("cron jobs") and asynchronous sweeps. This architectural decision avoids the need for OS-level cron configurations (like `crontab`) and ensures that all background tasks share the same connection pools, configuration, and environment as the main Core API.

## Running the Worker

### Development
To run the background worker locally during development, use the provided Makefile command:

```bash
make worker
```
This executes `go run cmd/worker/main.go`.

### Production
In production, the worker should be run as a separate, persistent process alongside the Core API.

**Docker / Kubernetes:**
If you are deploying via Docker, use the same container image as the API but override the command/entrypoint:

```dockerfile
# Start the worker
CMD ["/app/worker"]
```

**Systemd (Linux):**
If you are deploying directly to a Linux VM, create a `kobo-worker.service` file:

```ini
[Unit]
Description=Kobo Background Worker
After=network.target postgresql.service

[Service]
Type=simple
User=kobo
WorkingDirectory=/opt/kobo
ExecStart=/opt/kobo/worker
Restart=always
RestartSec=5
EnvironmentFile=/opt/kobo/.env

[Install]
WantedBy=multi-user.target
```

## Configured Jobs (Cron Substitutes)

The worker process (`cmd/worker/main.go`) uses Go's `time.Ticker` to schedule tasks. Current scheduled jobs include:

### 1. Reconciliation Sweep
- **Interval:** Every 30 minutes
- **Description:** Scans the ledger for unsettled transactions and matches them against provisioned virtual accounts.
- **Implementation:** `reconciliation.Sweeper`

### 2. Closure Sweep
- **Interval:** Runs alongside the Reconciliation sweep
- **Description:** Sweeps balances out of accounts that have been marked for closure, ensuring funds are routed to the fallback sweep destination.
- **Implementation:** `reconciliation.ClosureSweeper`

### 3. KYC Checks
- **Interval:** Every 1 hour
- **Description:** Placeholder for background KYC status verification for integrators.

### 4. Automated Billing & Invoicing
- **Interval:** Every 12 hours
- **Description:** Processes the automated monthly billing cycle. Retrieves pending or failed invoices, queries the integrator's saved Nomba payment token, and triggers a background charge via the `ChargeToken` Nomba API. Handles max-retry logic (up to 3 times) and suspends delinquent accounts.
- **Implementation:** `billing.InvoiceJob`

## Modifying Job Intervals

If you need to change the frequency of a job, modify the respective `time.NewTicker` duration in `cmd/worker/main.go`:

```go
billingTicker := time.NewTicker(12 * time.Hour) // Change this value to adjust the billing cron schedule
```

Note that any change to the intervals requires restarting the worker process to take effect.
