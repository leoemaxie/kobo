<script lang="ts">
  import CurrentPeriodCard from './CurrentPeriodCard.svelte';
  import PlanDetailsCard from './PlanDetailsCard.svelte';
  import InvoiceHistory from './InvoiceHistory.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import { toast } from '$lib/state/toast.svelte';

  const state = useConsoleState();
</script>

<svelte:head>
  <title>Billing — Kobo Console</title>
</svelte:head>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Page bar -->
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    padding-bottom: 20px; border-bottom: 1px solid var(--border-subtle);
  ">
    <div>
      <p style="
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 6px;
      ">Billing & Usage</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">plan:</span>
        <code style="
          font-family: monospace; font-size: 13px; color: #C8C8C8;
          background: var(--bg-sidebar); border: 1px solid #2a2a2a;
          border-radius: 4px; padding: 2px 8px;
        ">{state.billingOverview?.plan || 'pay_as_you_go'}</code>
        <span style="font-size: 11px; color: var(--text-muted);">·</span>
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">
          next_invoice: <span style="color: #C8C8C8;">{state.billingOverview?.nextInvoiceDate || '2026-11-01'}</span>
        </span>
      </div>
    </div>
    <button on:click={() => toast.info('Payment gateway integration pending.')} style="
      border: 1px solid #2a2a2a; border-radius: 6px;
      background: var(--bg-sidebar); padding: 6px 12px;
      font-size: 13px; font-weight: 600; color: var(--text-muted); cursor: pointer;
    ">
      Manage Payment Method
    </button>
  </div>

  <!-- Metrics + plan side-by-side -->
  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: var(--border-subtle); border: 1px solid var(--border-subtle); border-radius: 8px; overflow: hidden;">
    <CurrentPeriodCard />
    <PlanDetailsCard />
  </div>

  <InvoiceHistory />
</div>
