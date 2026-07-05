<script lang="ts">
  let { label, id, type = 'text', value = $bindable(), required = false, class: className = '', ...rest } = $props();
</script>

<div class="space-y-1.5 {className}">
  {#if label}
    <label for={id} class="block text-sm font-medium text-main">
      {label}
      {#if required}
        <span class="text-primary ml-0.5">*</span>
      {/if}
    </label>
  {/if}
  <input
    {id}
    {type}
    {required}
    bind:value
    class="block w-full rounded-radius-md border border-border bg-background px-3 py-2 text-main placeholder-subtle shadow-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary sm:text-sm transition-colors"
    {...rest}
  />
</div>
