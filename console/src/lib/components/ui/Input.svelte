<script lang="ts">
  /**
   * Input — design-system form input.
   *
   * Variants:
   *  default  — full-width standard input (used in forms)
   *  settings — used in settings / support forms (max-width 400px, sidebar bg)
   */
  let {
    label,
    id,
    type = 'text',
    value = $bindable(),
    required = false,
    variant = 'default',
    class: className = '',
    hint = '',
    ...rest
  }: {
    label?: string;
    id?: string;
    type?: string;
    value?: string;
    required?: boolean;
    variant?: 'default' | 'settings';
    class?: string;
    hint?: string;
    [key: string]: unknown;
  } = $props();

  const inputStyles: Record<string, string> = {
    default:  'block w-full rounded-[6px] border border-border bg-element px-3 py-2 text-sm text-main placeholder:text-subtle shadow-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors disabled:opacity-60 disabled:cursor-not-allowed',
    settings: 'block w-full max-w-[400px] rounded-[6px] border border-border bg-element px-3 py-2 text-sm text-main placeholder:text-subtle focus:border-primary focus:outline-none transition-colors disabled:opacity-60 disabled:cursor-not-allowed',
  };
</script>

<div class="space-y-1.5 {className}">
  {#if label}
    <label for={id} class="block text-xs font-semibold text-muted uppercase tracking-widest">
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
    class={inputStyles[variant]}
    {...rest}
  />

  {#if hint}
    <p class="text-[11px] text-subtle mt-1">{hint}</p>
  {/if}
</div>
