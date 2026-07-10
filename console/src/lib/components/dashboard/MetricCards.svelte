<script lang="ts">
  import { ArrowUpRight, ArrowDownRight, Minus } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const metrics = $derived(state.metrics);
</script>

<div
  class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-px bg-[var(--border-color)] border border-[var(--border-color)] rounded-xl overflow-hidden"
>
  {#each metrics as m}
    <div class="bg-[var(--bg-sidebar)] px-6 py-5">
      <p
        style="
        font-size: 11px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 14px;
      "
      >
        {m.label}
      </p>

      <div style="display: flex; align-items: baseline; gap: 10px; margin-bottom: 6px;">
        <span
          style="
          font-family: 'JetBrains Mono', monospace; font-size: 27px;
          font-weight: 700; color: var(--text-main); letter-spacing: -0.5px;
        ">{m.value}</span
        >
        <span
          style="
          display: flex; align-items: center; gap: 3px;
          font-size: 12px; font-weight: 600;
          color: {m.trend === 'up'
            ? 'var(--accent)'
            : m.trend === 'down'
              ? 'var(--error-color)'
              : 'var(--text-muted)'};
        "
        >
          {#if m.trend === 'up'}<ArrowUpRight size={11} />
          {:else if m.trend === 'down'}<ArrowDownRight size={11} />
          {:else}<Minus size={11} />
          {/if}
          {m.delta}
        </span>
      </div>

      <div
        style="
        width: 100%; height: 2px; background: var(--border-color);
        border-radius: 2px; overflow: hidden; margin-bottom: 8px;
      "
      >
        <div
          style="
          height: 100%; width: {m.bar}%; border-radius: 2px;
          background: {m.trend === 'up'
            ? 'var(--accent)'
            : m.trend === 'down'
              ? 'var(--error-color)'
              : 'var(--text-subtle)'};
          box-shadow: {m.trend === 'up' ? '0 0 6px var(--accent-glow)' : 'none'};
        "
        ></div>
      </div>

      <p style="font-size: 12px; color: var(--text-subtle); margin: 0;">{m.sub}</p>
    </div>
  {/each}
</div>
