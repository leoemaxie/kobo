<script lang="ts">
  /**
   * SectionLabel — the tiny uppercase monospace label used above tables
   * and sections ("STANDARD KEYS", "INVOICE HISTORY", "CURRENT PERIOD · …")
   */
  let {
    children,
    class: className = '',
  }: {
    children: import('svelte').Snippet;
    class?: string;
  } = $props();
</script>

<p class="text-[10px] font-bold uppercase tracking-widest text-subtle mb-2.5 {className}">
  {@render children()}
</p>
