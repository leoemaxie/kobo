<script lang="ts">
  let { label, id, type = 'text', value = $bindable(), required = false, class: className = '', ...rest } = $props();
</script>

<div class="space-y-1.5 {className}">
  {#if label}
    <label for={id} class="block text-sm font-medium text-paper">
      {label}
      {#if required}
        <span class="text-electric-lime ml-0.5">*</span>
      {/if}
    </label>
  {/if}
  <input
    {id}
    {type}
    {required}
    bind:value
    class="block w-full rounded-radius-md border border-iron bg-void-black px-3 py-2 text-paper placeholder-fog shadow-sm focus:border-steel focus:outline-none focus:ring-1 focus:ring-steel sm:text-sm transition-colors"
    {...rest}
  />
</div>
