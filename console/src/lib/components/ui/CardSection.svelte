<script lang="ts">
  /**
   * CardSection — the card-style container with a header band (title + subtitle)
   * and a padded body area. Used in WorkspaceInfoCard, TeamList, support page etc.
   *
   * Slots / Snippets:
   *  - header: content rendered in the header band (right side, after title/subtitle)
   *  - children: card body
   */
  let {
    title,
    subtitle,
    header,
    children,
    class: className = '',
    bodyClass = 'p-5',
    danger = false,
  }: {
    title: string;
    subtitle?: string;
    header?: import('svelte').Snippet;
    children: import('svelte').Snippet;
    class?: string;
    bodyClass?: string;
    danger?: boolean;
  } = $props();
</script>

<div
  class="bg-element border rounded-lg overflow-hidden {danger
    ? 'border-red-500/25'
    : 'border-border-subtle'} {className}"
>
  <!-- Header band -->
  <div
    class="px-5 py-4 border-b flex items-center justify-between {danger
      ? 'border-red-500/25 bg-red-500/[0.06]'
      : 'border-border-subtle bg-sidebar'}"
  >
    <div>
      <h3 class="text-sm font-semibold {danger ? 'text-red-400' : 'text-main'} m-0">{title}</h3>
      {#if subtitle}
        <p class="text-[12px] text-subtle mt-1 mb-0">{subtitle}</p>
      {/if}
    </div>
    {#if header}
      <div>
        {@render header()}
      </div>
    {/if}
  </div>

  <!-- Body -->
  <div class={bodyClass}>
    {@render children()}
  </div>
</div>
