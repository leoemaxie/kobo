<script lang="ts">
  import CurrentPeriodCard from './CurrentPeriodCard.svelte';
  import PlanDetailsCard from './PlanDetailsCard.svelte';
  import InvoiceHistory from './InvoiceHistory.svelte';
</script>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Page bar -->
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    padding-bottom: 20px; border-bottom: 1px solid #1e1e1e;
  ">
    <div>
      <p style="
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: #555; margin: 0 0 6px;
      ">Billing & Usage</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: #555;">plan:</span>
        <code style="
          font-family: monospace; font-size: 13px; color: #C8C8C8;
          background: #111; border: 1px solid #2a2a2a;
          border-radius: 4px; padding: 2px 8px;
        ">pay_as_you_go</code>
        <span style="font-size: 11px; color: #444;">·</span>
        <span style="font-family: monospace; font-size: 13px; color: #555;">
          next_invoice: <span style="color: #C8C8C8;">2026-11-01</span>
        </span>
      </div>
    </div>
    <button style="
      border: 1px solid #2a2a2a; border-radius: 6px;
      background: #111; padding: 6px 12px;
      font-size: 13px; font-weight: 600; color: #888; cursor: pointer;
    ">
      Manage Payment Method
    </button>
  </div>

  <!-- Metrics + plan side-by-side -->
  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: #1e1e1e; border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
    <CurrentPeriodCard />
    <PlanDetailsCard />
  </div>

  <InvoiceHistory />
</div>
