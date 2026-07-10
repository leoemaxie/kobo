<script lang="ts">
  import { ArrowUpRight } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const rates = $derived(state.billingOverview?.planDetails || []);
</script>

<div
  style="background: var(--bg-element); padding: 24px; border-left: 1px solid var(--border-subtle);"
>
  <div
    style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;"
  >
    <p
      style="
      font-size: 11px; font-weight: 700; text-transform: uppercase;
      letter-spacing: 0.1em; color: var(--text-muted); margin: 0;
    "
    >
      Plan · pay_as_you_go
    </p>
    <a
      href="/billing/upgrade"
      style="
      display: flex; align-items: center; gap: 4px; text-decoration: none;
      font-size: 13px; font-weight: 600; color: var(--accent);
    "
    >
      Upgrade <ArrowUpRight size={11} />
    </a>
  </div>

  <div style="border: 1px solid var(--border-subtle); border-radius: 6px; overflow: hidden;">
    {#each rates as r, i}
      <div
        style="
        display: grid; grid-template-columns: 1fr auto;
        padding: 9px 12px;
        border-bottom: {i < rates.length - 1 ? '1px solid var(--bg-sidebar)' : 'none'};
        align-items: center;
      "
      >
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);"
          >{r.key}</span
        >
        <div style="text-align: right;">
          <span style="font-family: monospace; font-size: 14px; font-weight: 700; color: #C8C8C8;"
            >{r.unit}</span
          >
          <span
            style="font-family: monospace; font-size: 12px; color: var(--text-muted); margin-left: 5px;"
            >{r.label}</span
          >
        </div>
      </div>
    {/each}
  </div>

  <p
    style="font-family: monospace; font-size: 12px; color: var(--text-muted); margin: 12px 0 0; line-height: 1.6;"
  >
    No commitments. Usage billed monthly in arrears.
  </p>
</div>
