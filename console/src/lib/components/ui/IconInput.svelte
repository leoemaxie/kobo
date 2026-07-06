<script lang="ts">
  /**
   * IconInput — styled input with an optional leading icon and trailing toggle
   * Used on all auth pages (email + password fields).
   *
   * Props mirror the standard <input> element, plus:
   *  - label: optional string label
   *  - icon: optional Snippet rendered on the left
   *  - trailing: optional Snippet rendered on the right (e.g. eye toggle)
   *  - required: boolean
   */
  import type { Snippet } from 'svelte';

  let {
    id,
    label = '',
    type = 'text',
    placeholder = '',
    required = false,
    value = $bindable(''),
    name = '',
    icon,
    trailing,
    class: className = '',
    labelClass = '',
    ...rest
  }: {
    id?: string;
    label?: string;
    type?: string;
    placeholder?: string;
    required?: boolean;
    value?: string;
    name?: string;
    icon?: Snippet;
    trailing?: Snippet;
    class?: string;
    labelClass?: string;
    [key: string]: unknown;
  } = $props();
</script>

<div class="space-y-1.5 {className}">
  {#if label}
    <label for={id} class="block text-xs font-semibold text-muted uppercase tracking-widest {labelClass}">
      {label}
      {#if required}<span class="text-primary ml-0.5">*</span>{/if}
    </label>
  {/if}

  <div class="relative">
    {#if icon}
      <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
        {@render icon()}
      </div>
    {/if}

    <input
      {id}
      {name}
      {type}
      {required}
      {placeholder}
      bind:value
      class="block w-full rounded-[8px] border border-border bg-background
             {icon ? 'pl-10' : 'pl-4'}
             {trailing ? 'pr-10' : 'pr-4'}
             py-3 text-sm text-main placeholder-subtle
             focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary
             transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
      {...rest}
    />

    {#if trailing}
      <div class="absolute inset-y-0 right-0 flex items-center">
        {@render trailing()}
      </div>
    {/if}
  </div>
</div>
