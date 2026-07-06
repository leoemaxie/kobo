<script lang="ts">
  import { toast } from '$lib/state/toast.svelte';
  import { X, CheckCircle, AlertTriangle, Info } from '@lucide/svelte';
  import { fly } from 'svelte/transition';

  const toasts = $derived(toast.toasts);

  // Track remaining % for each toast's progress bar
  let progress = $state<Record<string, number>>({});

  $effect(() => {
    for (const t of toasts) {
      if (!(t.id in progress)) {
        progress[t.id] = 100;
        const interval = setInterval(() => {
          const elapsed = Date.now() - t.createdAt;
          const remaining = Math.max(0, 100 - (elapsed / t.duration) * 100);
          progress[t.id] = remaining;
          if (remaining === 0) clearInterval(interval);
        }, 30);
      }
    }
  });

  const ACCENT_MAP = {
    success: 'var(--accent)',
    error: '#f87171',
    info: '#60a5fa',
  };

  const LABEL_MAP = {
    success: 'Success',
    error: 'Error',
    info: 'Info',
  };
</script>

<div
  class="fixed bottom-5 right-5 z-[9999] flex flex-col-reverse gap-2.5 w-full max-w-[360px] pointer-events-none"
  aria-live="polite"
  aria-label="Notifications"
>
  {#each toasts as t (t.id)}
    <div
      in:fly={{ y: 16, duration: 280, opacity: 0 }}
      out:fly={{ y: 8, duration: 200, opacity: 0 }}
      class="pointer-events-auto relative flex flex-col rounded-[10px] overflow-hidden border border-border"
      style="
        background: var(--bg-element);
        box-shadow: 0 4px 24px rgba(0,0,0,0.18), 0 1px 4px rgba(0,0,0,0.1);
      "
    >
      <!-- Top content row -->
      <div class="flex items-start gap-2.5 px-3.5 pt-3 pb-2.5">

        <!-- Icon -->
        <div class="flex-shrink-0 mt-[1px]">
          {#if t.type === 'success'}
            <CheckCircle size={14} style="color: {ACCENT_MAP.success};" />
          {:else if t.type === 'error'}
            <AlertTriangle size={14} style="color: {ACCENT_MAP.error};" />
          {:else}
            <Info size={14} style="color: {ACCENT_MAP.info};" />
          {/if}
        </div>

        <!-- Text -->
        <div class="flex-1 min-w-0">
          <p
            class="text-[11px] font-bold uppercase tracking-widest mb-0.5"
            style="color: {ACCENT_MAP[t.type]}; letter-spacing: 0.08em;"
          >
            {LABEL_MAP[t.type]}
          </p>
          <p class="text-[12.5px] font-medium leading-snug text-main m-0">
            {t.message}
          </p>
        </div>

        <!-- Dismiss -->
        <button
          onclick={() => toast.remove(t.id)}
          class="flex-shrink-0 flex items-center justify-center w-5 h-5 rounded-md cursor-pointer transition-colors text-subtle hover:text-main"
          style="background: transparent; border: none;"
          onmouseenter={(e) =>
            ((e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-active)')}
          onmouseleave={(e) =>
            ((e.currentTarget as HTMLButtonElement).style.background = 'transparent')}
          aria-label="Dismiss"
        >
          <X size={12} />
        </button>
      </div>

      <!-- Progress bar -->
      <div
        class="h-[2px] w-full"
        style="background: var(--bg-active);"
      >
        <div
          class="h-full transition-none"
          style="
            width: {progress[t.id] ?? 100}%;
            background: {ACCENT_MAP[t.type]};
            opacity: 0.7;
          "
        ></div>
      </div>
    </div>
  {/each}
</div>
