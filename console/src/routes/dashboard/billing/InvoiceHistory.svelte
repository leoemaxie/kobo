<script lang="ts">
  import { Download } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const invoices = $derived(state.billingInvoices);

  const cols = ['INVOICE ID', 'PERIOD', 'AMOUNT', 'STATUS', 'DATE', ''];
</script>

<div>
  <p style="
    font-size: 12px; font-weight: 700; text-transform: uppercase;
    letter-spacing: 0.1em; color: #555; margin: 0 0 10px;
  ">Invoice History</p>

  <div style="border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
    <!-- Header -->
    <div style="
      display: grid; grid-template-columns: 1.4fr 1fr 1fr 80px 100px 40px;
      padding: 8px 16px; background: #0d0d0d; border-bottom: 1px solid #1e1e1e;
    ">
      {#each cols as col}
        <span style="font-size: 11px; font-weight: 700; letter-spacing: 0.1em; color: #444; text-transform: uppercase;">
          {col}
        </span>
      {/each}
    </div>

    {#each invoices as inv}
      <div style="
        display: grid; grid-template-columns: 1.4fr 1fr 1fr 80px 100px 40px;
        padding: 11px 16px; align-items: center; border-bottom: 1px solid #111;
        transition: background 0.1s;
      "
        onmouseenter={(e) => (e.currentTarget as HTMLDivElement).style.background = '#0f0f0f'}
        onmouseleave={(e) => (e.currentTarget as HTMLDivElement).style.background = 'transparent'}
      >
        <code style="font-family: monospace; font-size: 13px; color: #666;">{inv.id}</code>
        <span style="font-size: 14px; color: #C8C8C8;">{inv.period}</span>
        <span style="font-family: monospace; font-size: 14px; font-weight: 600; color: #C8C8C8;">{inv.amount}</span>
        <span style="
          font-family: monospace; font-size: 12px; font-weight: 700; text-transform: uppercase;
          color: #C0FF00;
        ">{inv.status}</span>
        <span style="font-family: monospace; font-size: 13px; color: #555;">{inv.date}</span>
        <button style="
          background: none; border: none; cursor: pointer; color: #555;
          display: flex; padding: 0;
        " title="Download PDF">
          <Download size={13} />
        </button>
      </div>
    {/each}
  </div>
</div>
