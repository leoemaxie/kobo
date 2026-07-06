<script lang="ts">
  /**
   * Reusable dashboard page header bar.
   *
   * Props:
   *  - title: string  — the main UPPERCASE label (e.g. "API Keys")
   *  - meta: optional Snippet — secondary content rendered after the title (code badges, env badge, etc.)
   *  - actions: optional Snippet — rendered on the right side (e.g. CTA buttons)
   */
  let {
    title,
    meta,
    actions,
  }: {
    title: string;
    meta?: import('svelte').Snippet;
    actions?: import('svelte').Snippet;
  } = $props();
</script>

<div class="flex items-center justify-between pb-5 border-b border-border-subtle">
  <div>
    <p class="text-[11px] font-bold uppercase tracking-widest text-subtle mb-1.5">{title}</p>
    {#if meta}
      <div class="flex items-center gap-2">
        {@render meta()}
      </div>
    {/if}
  </div>
  {#if actions}
    <div class="flex items-center gap-2">
      {@render actions()}
    </div>
  {/if}
</div>
