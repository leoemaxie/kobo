<script lang="ts">
  import { ChevronDown, Check } from '@lucide/svelte';
  import { onMount, onDestroy } from 'svelte';

  export interface SelectOption {
    value: string;
    label: string;
  }

  let {
    options = [],
    value = $bindable(''),
    placeholder = 'Select an option',
    name = '',
    id = '',
    required = false,
  } = $props<{
    options: SelectOption[];
    value?: string;
    placeholder?: string;
    name?: string;
    id?: string;
    required?: boolean;
  }>();

  let isOpen = $state(false);
  let containerRef: HTMLDivElement;

  function toggleOpen() {
    isOpen = !isOpen;
  }

  function selectOption(val: string) {
    value = val;
    isOpen = false;
  }

  function handleClickOutside(e: MouseEvent) {
    if (containerRef && !containerRef.contains(e.target as Node)) {
      isOpen = false;
    }
  }

  onMount(() => {
    document.addEventListener('click', handleClickOutside);
  });
  
  onDestroy(() => {
    if (typeof document !== 'undefined') {
      document.removeEventListener('click', handleClickOutside);
    }
  });

  let selectedLabel = $derived(options.find((o: SelectOption) => o.value === value)?.label || placeholder);
</script>

<div class="relative w-full {isOpen ? 'z-50' : ''}" bind:this={containerRef}>
  {#if name}
    <input type="hidden" {name} {value} {required} {id} />
  {/if}

  <button
    type="button"
    class="w-full h-9 px-3 rounded-lg border border-[var(--border-color)] bg-[var(--bg-element)] text-[13px] outline-none focus:border-[var(--accent)] flex items-center justify-between transition-colors duration-200"
    onclick={toggleOpen}
    class:border-[var(--accent)]={isOpen}
  >
    <span class="truncate {value ? 'text-main' : 'text-muted'}">{selectedLabel}</span>
    <ChevronDown size={14} class="text-muted transition-transform duration-200 {isOpen ? 'rotate-180' : ''}" />
  </button>

  {#if isOpen}
    <div
      class="absolute z-[100] top-full left-0 right-0 mt-1.5 max-h-52 overflow-y-auto rounded-lg border border-[var(--border-color)] bg-[var(--bg-element)] shadow-xl p-1 custom-scrollbar"
      style="animation: select-slide-down 0.15s ease-out forwards;"
    >
      {#each options as option}
        <button
          type="button"
          class="w-full text-left px-2.5 py-1.5 rounded-md text-[13px] flex items-center justify-between transition-colors duration-150 {value === option.value ? 'bg-[var(--accent)]/10 text-[var(--accent)] font-medium' : 'text-main hover:bg-[var(--bg-element)]'}"
          onclick={() => selectOption(option.value)}
        >
          <span class="truncate">{option.label}</span>
          {#if value === option.value}
            <Check size={14} class="text-[var(--accent)] flex-shrink-0 ml-2" />
          {/if}
        </button>
      {/each}
      
      {#if options.length === 0}
        <div class="px-2.5 py-2 text-[13px] text-muted text-center">No options available</div>
      {/if}
    </div>
  {/if}
</div>

<style>
  @keyframes select-slide-down {
    from {
      opacity: 0;
      transform: translateY(-4px) scale(0.98);
    }
    to {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }

  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: var(--border-color);
    border-radius: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: var(--text-muted);
  }
</style>
