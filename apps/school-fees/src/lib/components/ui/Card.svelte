<script lang="ts">
  let { children, class: className = '', title = '', ...rest } = $props();
</script>

<div
  class="bg-carbon border border-iron rounded-radius-cards p-6 sm:p-8 shadow-sm {className}"
  {...rest}
>
  {#if title}
    <h3 class="font-basier text-xl font-semibold text-pure-white mb-4">{title}</h3>
  {/if}
  {@render children()}
</div>
