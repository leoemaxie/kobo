<script lang="ts">
  import { X } from '@lucide/svelte';

  let {
    title,
    onClose,
    children,
    maxWidth = 'max-w-md',
    scrollable = false,
    class: className = '',
  }: {
    title: string;
    onClose: () => void;
    children: import('svelte').Snippet;
    maxWidth?: string;
    scrollable?: boolean;
    class?: string;
  } = $props();
</script>

<div
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm {scrollable
    ? 'overflow-y-auto'
    : ''}"
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
>
  <div
    class="w-full {maxWidth} {scrollable
      ? 'my-8'
      : ''} bg-element border border-border rounded-xl overflow-hidden shadow-xl {className}"
  >
    <!-- Header -->
    <div class="flex items-center justify-between px-6 py-4 border-b border-border bg-sidebar">
      <h2 id="modal-title" class="text-sm font-semibold text-main uppercase tracking-widest">
        {title}
      </h2>
      <button
        onclick={onClose}
        class="text-muted hover:text-main transition-colors"
        aria-label="Close modal"
      >
        <X size={16} />
      </button>
    </div>

    <!-- Body -->
    {@render children()}
  </div>
</div>
