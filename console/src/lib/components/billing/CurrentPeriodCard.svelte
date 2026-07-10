<script lang="ts">
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const usageItems = $derived(state.billingOverview?.usageItems || []);

</script>

<div style="background: var(--bg-element); padding: 24px;">
  <p style="
    font-size: 11px; font-weight: 700; text-transform: uppercase;
    letter-spacing: 0.1em; color: var(--text-muted); margin: 0 0 16px;
  ">Current Period · {state.billingOverview?.period}</p>

  <div style="display: flex; align-items: baseline; gap: 10px; margin-bottom: 20px;">
    <span style="
      font-family: 'JetBrains Mono', monospace; font-size: 34px;
      font-weight: 700; color: var(--text-main); letter-spacing: -1px;
    ">{state.billingOverview?.accrued}</span>
    <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">accrued</span>
  </div>

  <div style="display: flex; flex-direction: column; gap: 14px;">
    {#each usageItems as item}
      <div>
        <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 5px;">
          <span style="font-family: monospace; font-size: 12px; color: var(--text-muted);">{item.key}</span>
          <span style="font-family: monospace; font-size: 12px; color: var(--text-subtle);">
            {item.calc} = <span style="color: #C8C8C8;">{item.amount}</span>
          </span>
        </div>
        <div style="
          width: 100%; height: 2px; background: var(--bg-active); border-radius: 2px; overflow: hidden;
        ">
          <div style="
            height: 100%; width: {item.pct}%; background: var(--accent); border-radius: 2px;
            box-shadow: 0 0 5px var(--accent-glow);
          "></div>
        </div>
      </div>
    {/each}
  </div>
</div>
