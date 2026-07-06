<script lang="ts">
  /**
   * AuthLogo — the brand logo header block used at the top of every auth page.
   * Renders the Kobo logo, a heading, and optional subtitle/link.
   */
  let {
    heading,
    children,
  }: {
    heading: string;
    children?: import('svelte').Snippet;
  } = $props();
</script>

<div class="mb-8 text-center">
  <img
    src="/logo.png"
    alt="Kobo"
    class="h-10 w-auto mx-auto mb-6"
    style="filter: var(--logo-filter);"
  />
  <h1 class="text-2xl font-inter font-bold text-main tracking-tight">{heading}</h1>
  {#if children}
    <div class="text-muted text-sm mt-2 leading-relaxed">
      {@render children()}
    </div>
  {/if}
</div>
