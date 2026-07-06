<script lang="ts">
  let {
    label,
    id,
    rows = 4,
    value = $bindable(''),
    required = false,
    placeholder = '',
    name = '',
    variant = 'default',
    class: className = '',
    ...rest
  }: {
    label?: string;
    id?: string;
    rows?: number;
    value?: string;
    required?: boolean;
    placeholder?: string;
    name?: string;
    variant?: 'default' | 'settings' | 'mono';
    class?: string;
    [key: string]: unknown;
  } = $props();

  const inputStyles: Record<string, string> = {
    default:  'block w-full rounded-[6px] border border-border bg-sidebar px-3 py-2 text-sm text-main placeholder:text-subtle focus:border-primary focus:outline-none transition-colors resize-vertical',
    settings: 'block w-full max-w-[400px] rounded-[6px] border border-border bg-element px-3 py-2 text-sm text-main placeholder:text-subtle focus:border-primary focus:outline-none transition-colors resize-vertical',
    mono:     'block w-full rounded-lg border border-border bg-sidebar px-4 py-2.5 text-sm font-mono text-main placeholder:text-subtle focus:border-border-focus focus:outline-none transition-colors resize-vertical',
  };
</script>

<div class="space-y-1.5 {className}">
  {#if label}
    <label for={id} class="block text-xs font-semibold text-muted uppercase tracking-widest">
      {label}
      {#if required}<span class="text-primary ml-0.5">*</span>{/if}
    </label>
  {/if}

  <textarea
    {id}
    {name}
    {rows}
    {required}
    {placeholder}
    bind:value
    class={inputStyles[variant]}
    {...rest}
  ></textarea>
</div>
