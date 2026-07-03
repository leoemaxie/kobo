<script lang="ts">
  import { ArrowUpRight } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const rates = $derived(state.billingOverview?.planDetails || []);
</script>

<div style="background: #0a0a0a; padding: 24px; border-left: 1px solid #1e1e1e;">
  <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
    <p style="
      font-size: 11px; font-weight: 700; text-transform: uppercase;
      letter-spacing: 0.1em; color: #444; margin: 0;
    ">Plan · pay_as_you_go</p>
    <a href="/billing/upgrade" style="
      display: flex; align-items: center; gap: 4px; text-decoration: none;
      font-size: 13px; font-weight: 600; color: #C0FF00;
    ">
      Upgrade <ArrowUpRight size={11} />
    </a>
  </div>

  <div style="border: 1px solid #1e1e1e; border-radius: 6px; overflow: hidden;">
    {#each rates as r, i}
      <div style="
        display: grid; grid-template-columns: 1fr auto;
        padding: 9px 12px;
        border-bottom: {i < rates.length - 1 ? '1px solid #111' : 'none'};
        align-items: center;
      ">
        <span style="font-family: monospace; font-size: 13px; color: #555;">{r.key}</span>
        <div style="text-align: right;">
          <span style="font-family: monospace; font-size: 14px; font-weight: 700; color: #C8C8C8;">{r.unit}</span>
          <span style="font-family: monospace; font-size: 12px; color: #444; margin-left: 5px;">{r.label}</span>
        </div>
      </div>
    {/each}
  </div>

  <p style="font-family: monospace; font-size: 12px; color: #444; margin: 12px 0 0; line-height: 1.6;">
    No commitments. Usage billed monthly in arrears.
  </p>
</div>
