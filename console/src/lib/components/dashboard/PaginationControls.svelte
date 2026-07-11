<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { ChevronLeft, ChevronRight } from '@lucide/svelte';

  let { meta } = $props<{ meta: { total: number; page: number; limit: number; totalPages: number } }>();

  function setPage(p: number) {
    if (p < 1 || p > meta.totalPages || p === meta.page) return;
    const url = new URL(page.url);
    url.searchParams.set('page', p.toString());
    goto(url.toString(), { keepFocus: true, noScroll: false });
  }

  let startItem = $derived((meta.page - 1) * meta.limit + 1);
  let endItem = $derived(Math.min(meta.page * meta.limit, meta.total));
</script>

{#if meta.total > 0}
  <div class="flex items-center justify-between mt-6">
    <p class="text-[13px] text-[var(--text-subtle)]">
      Showing <span class="font-semibold text-[var(--text-main)]">{startItem}</span> to
      <span class="font-semibold text-[var(--text-main)]">{endItem}</span> of
      <span class="font-semibold text-[var(--text-main)]">{meta.total}</span> entries
    </p>

    <div class="flex items-center gap-2">
      <button
        onclick={() => setPage(meta.page - 1)}
        disabled={meta.page === 1}
        class="flex items-center justify-center w-8 h-8 rounded border border-[var(--border-color)] bg-transparent text-[var(--text-main)] hover:bg-[var(--bg-element)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        aria-label="Previous page"
      >
        <ChevronLeft size={16} />
      </button>

      <span class="text-[13px] text-[var(--text-main)] mx-2">
        Page <span class="font-semibold">{meta.page}</span> of {meta.totalPages}
      </span>

      <button
        onclick={() => setPage(meta.page + 1)}
        disabled={meta.page === meta.totalPages}
        class="flex items-center justify-center w-8 h-8 rounded border border-[var(--border-color)] bg-transparent text-[var(--text-main)] hover:bg-[var(--bg-element)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        aria-label="Next page"
      >
        <ChevronRight size={16} />
      </button>
    </div>
  </div>
{/if}
