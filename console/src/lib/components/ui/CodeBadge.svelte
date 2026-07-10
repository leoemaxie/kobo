<script lang="ts">
  /**
   * Accent code badge — the monospace highlighted pill used to display
   * environment names, IDs, plan names etc. throughout the console.
   *
   * Props:
   *  - variant: 'accent' | 'neutral'
   *    accent  → green accent background / border (default)
   *    neutral → sidebar background / border
   */
  let {
    children,
    variant = 'accent',
    class: className = '',
  }: {
    children: import('svelte').Snippet;
    variant?: 'accent' | 'neutral';
    class?: string;
  } = $props();

  const styles: Record<string, string> = {
    accent: 'bg-primary-transparent border-primary-border text-primary',
    neutral: 'bg-active border-border text-muted',
  };
</script>

<code
  class="font-inconsolata text-[11px] border rounded px-2 py-0.5 tracking-[0.05em] {styles[
    variant
  ]} {className}"
>
  {@render children()}
</code>
