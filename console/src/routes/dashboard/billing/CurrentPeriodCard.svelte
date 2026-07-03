<script lang="ts">
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const usageItems = $derived(state.billingOverview?.usageItems || []);

</script>

<div style="background: #0a0a0a; padding: 24px;">
  <p style="
    font-size: 11px; font-weight: 700; text-transform: uppercase;
    letter-spacing: 0.1em; color: #444; margin: 0 0 16px;
  ">Current Period · {state.billingOverview?.period}</p>

  <div style="display: flex; align-items: baseline; gap: 10px; margin-bottom: 20px;">
    <span style="
      font-family: 'JetBrains Mono', monospace; font-size: 34px;
      font-weight: 700; color: #F8F8F8; letter-spacing: -1px;
    ">{state.billingOverview?.accrued}</span>
    <span style="font-family: monospace; font-size: 13px; color: #555;">accrued</span>
  </div>

  <div style="display: flex; flex-direction: column; gap: 14px;">
    {#each usageItems as item}
      <div>
        <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 5px;">
          <span style="font-family: monospace; font-size: 12px; color: #888;">{item.key}</span>
          <span style="font-family: monospace; font-size: 12px; color: #555;">
            {item.calc} = <span style="color: #C8C8C8;">{item.amount}</span>
          </span>
        </div>
        <div style="
          width: 100%; height: 2px; background: #1a1a1a; border-radius: 2px; overflow: hidden;
        ">
          <div style="
            height: 100%; width: {item.pct}%; background: #C0FF00; border-radius: 2px;
            box-shadow: 0 0 5px rgba(192,255,0,0.4);
          "></div>
        </div>
      </div>
    {/each}
  </div>
</div>
