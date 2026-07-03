<script lang="ts">
  import { ArrowUpRight, ArrowDownRight, Minus } from '@lucide/svelte';

  const metrics = [
    {
      key: 'api_requests',
      label: 'API Requests',
      value: '1,248',
      delta: '+12.5%',
      sub: 'Last 30 days',
      trend: 'up',
      bar: 65,
    },
    {
      key: 'virtual_accounts',
      label: 'Virtual Accounts',
      value: '342',
      delta: '+4',
      sub: '500 limit (sandbox)',
      trend: 'up',
      bar: 68,
    },
    {
      key: 'error_rate',
      label: 'Error Rate',
      value: '0.8%',
      delta: '−0.2pp',
      sub: 'vs. prior period',
      trend: 'down',
      bar: 8,
    },
    {
      key: 'p99_latency',
      label: 'p99 Latency',
      value: '142ms',
      delta: '~',
      sub: 'stable this week',
      trend: 'neutral',
      bar: 28,
    },
  ];
</script>

<div style="
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 1px;
  border: 1px solid #1e1e1e; border-radius: 10px; overflow: hidden; background: #1e1e1e;
">
  {#each metrics as m}
    <div style="background: #0a0a0a; padding: 20px 24px;">
      <p style="
        font-size: 11px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: #555; margin: 0 0 14px;
      ">{m.label}</p>

      <div style="display: flex; align-items: baseline; gap: 10px; margin-bottom: 6px;">
        <span style="
          font-family: 'JetBrains Mono', monospace; font-size: 27px;
          font-weight: 700; color: #F8F8F8; letter-spacing: -0.5px;
        ">{m.value}</span>
        <span style="
          display: flex; align-items: center; gap: 3px;
          font-size: 12px; font-weight: 600;
          color: {m.trend === 'up' ? '#C0FF00' : m.trend === 'down' ? '#f87171' : '#666'};
        ">
          {#if m.trend === 'up'}<ArrowUpRight size={11} />
          {:else if m.trend === 'down'}<ArrowDownRight size={11} />
          {:else}<Minus size={11} />
          {/if}
          {m.delta}
        </span>
      </div>

      <div style="
        width: 100%; height: 2px; background: #1a1a1a;
        border-radius: 2px; overflow: hidden; margin-bottom: 8px;
      ">
        <div style="
          height: 100%; width: {m.bar}%; border-radius: 2px;
          background: {m.trend === 'up' ? '#C0FF00' : m.trend === 'down' ? '#f87171' : '#444'};
          box-shadow: {m.trend === 'up' ? '0 0 6px rgba(192,255,0,0.5)' : 'none'};
        "></div>
      </div>

      <p style="font-size: 12px; color: #555; margin: 0;">{m.sub}</p>
    </div>
  {/each}
</div>
