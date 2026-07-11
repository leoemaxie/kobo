<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    subtitle,
    rightSection,
  }: {
    title: string;
    subtitle?: string;
    rightSection?: Snippet;
  } = $props();
</script>

<div class="flex items-center justify-between border-b border-iron pb-6">
  <div>
    <h1 class="text-3xl font-bold text-pure-white tracking-tight">{title}</h1>
    {#if subtitle}
      <p class="text-smoke mt-1">{subtitle}</p>
    {/if}
  </div>
  {#if rightSection}
    {@render rightSection()}
  {/if}
</div>
