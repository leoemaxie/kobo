<script lang="ts">
  import { Download } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const consoleState = useConsoleState();
  const invoices = $derived(consoleState.billingInvoices);

  const cols = ['INVOICE ID', 'PERIOD', 'AMOUNT', 'STATUS', 'DATE', ''];
</script>

<div>
  <p style="
    font-size: 12px; font-weight: 700; text-transform: uppercase;
    letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 10px;
  ">Invoice History</p>

  <div style="border: 1px solid var(--border-subtle); border-radius: 8px; overflow: hidden;">
    <!-- Header -->
    <div style="
      display: grid; grid-template-columns: 1.4fr 1fr 1fr 80px 100px 40px;
      padding: 8px 16px; background: var(--bg-sidebar); border-bottom: 1px solid var(--border-subtle);
    ">
      {#each cols as col}
        <span style="font-size: 11px; font-weight: 700; letter-spacing: 0.1em; color: var(--text-muted); text-transform: uppercase;">
          {col}
        </span>
      {/each}
    </div>

    {#each invoices as inv}
      <div role="row" style="
        display: grid; grid-template-columns: 1.4fr 1fr 1fr 80px 100px 40px;
        padding: 11px 16px; align-items: center; border-bottom: 1px solid var(--bg-sidebar);
        transition: background 0.1s;
      "
        onmouseenter={(e) => (e.currentTarget as HTMLDivElement).style.background = '#0f0f0f'}
        onmouseleave={(e) => (e.currentTarget as HTMLDivElement).style.background = 'transparent'}
      >
        <code style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">{inv.id}</code>
        <span style="font-size: 14px; color: #C8C8C8;">{inv.period}</span>
        <span style="font-family: monospace; font-size: 14px; font-weight: 600; color: #C8C8C8;">{inv.amount}</span>
        <span style="
          font-family: monospace; font-size: 12px; font-weight: 700; text-transform: uppercase;
          color: var(--accent);
        ">{inv.status}</span>
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">{inv.date}</span>
        <button style="
          background: none; border: none; cursor: pointer; color: var(--text-subtle);
          display: flex; padding: 0;
        " title="Download PDF">
          <Download size={13} />
        </button>
      </div>
    {/each}
  </div>
</div>
