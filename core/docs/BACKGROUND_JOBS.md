# Background Worker and Scheduled Jobs

The Kobo Core platform utilizes a Go background worker process to handle scheduled tasks and asynchronous sweeps. This worker is designed as a **one-off execution script** intended to be triggered periodically via a cron scheduler, such as a GitHub Actions workflow or a Kubernetes CronJob.

This architectural decision ensures that background tasks are executed predictably, without the risk of long-running daemon memory leaks, while still sharing the same connection pools, configuration, and environment parsing as the main Core API.

## Running the Worker

### Development
To run the background worker locally during development for a single execution pass, use the provided Makefile command:

```bash
make worker
```
This executes `go run cmd/worker/main.go`.

### Production (GitHub Actions)
In production, the worker is configured to run via a GitHub Actions scheduled workflow. 

**Workflow File:** `.github/workflows/run-worker.yml`

This workflow checks out the code, sets up Go, and executes the worker using injected environment variables from your GitHub repository secrets. The cron schedule is defined within this file.

### Production (Docker / Kubernetes)
If you prefer deploying via Kubernetes, you can utilize a `CronJob` resource using the same container image as the API, overriding the command:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: kobo-worker
spec:
  schedule: "0 * * * *" # Every hour
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: worker
            image: kobo-core:latest
            command: ["/app/worker"]
            envFrom:
            - secretRef:
                name: kobo-secrets
          restartPolicy: OnFailure
```

## Configured Jobs

The worker process (`cmd/worker/main.go`) executes the following tasks sequentially on every run:

### 1. Reconciliation Sweep
- **Description:** Scans the ledger for unsettled transactions and matches them against provisioned virtual accounts.
- **Implementation:** `reconciliation.Sweeper`

### 2. Closure Sweep
- **Description:** Sweeps balances out of accounts that have been marked for closure, ensuring funds are routed to the fallback sweep destination.
- **Implementation:** `reconciliation.ClosureSweeper`

### 3. KYC Checks
- **Description:** Placeholder for background KYC status verification for integrators.

### 4. Automated Billing & Invoicing
- **Description:** Processes the automated billing cycle. Retrieves pending or failed invoices, queries the integrator's saved Monnify payment token, and triggers a background charge via the `ChargeToken` Monnify API. Handles max-retry logic and suspends delinquent accounts.
- **Implementation:** `billing.InvoiceJob`

## Modifying Job Execution
Because the worker is a one-off execution, the frequency of these jobs is determined entirely by the external scheduler (e.g., the GitHub Actions cron string). To change the frequency, update the `schedule` block in `.github/workflows/run-worker.yml`.
