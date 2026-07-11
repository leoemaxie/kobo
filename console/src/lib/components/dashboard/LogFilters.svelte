<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';

  let method = $state(page.url.searchParams.get('method') || '');
  let statusCode = $state(page.url.searchParams.get('status_code') || '');

  $effect(() => {
    method = page.url.searchParams.get('method') || '';
    statusCode = page.url.searchParams.get('status_code') || '';
  });

  function applyFilters() {
    const url = new URL(page.url);
    if (method) url.searchParams.set('method', method);
    else url.searchParams.delete('method');

    if (statusCode) url.searchParams.set('status_code', statusCode);
    else url.searchParams.delete('status_code');

    url.searchParams.set('page', '1'); // Reset to first page
    goto(url.toString(), { keepFocus: true });
  }

  function clearFilters() {
    method = '';
    statusCode = '';
    applyFilters();
  }
</script>

<div class="flex flex-wrap gap-4 mb-6 items-end">
  <div class="flex flex-col gap-1.5">
    <label
      for="method-filter"
      class="text-xs font-semibold text-[var(--text-subtle)] uppercase tracking-wider">Method</label
    >
    <select
      id="method-filter"
      bind:value={method}
      class="h-9 px-3 rounded-lg border border-[var(--border-color)] bg-transparent text-[13px] text-[var(--text-main)] outline-none focus:border-[var(--accent)] cursor-pointer"
    >
      <option value="">All Methods</option>
      <option value="GET">GET</option>
      <option value="POST">POST</option>
      <option value="PUT">PUT</option>
      <option value="PATCH">PATCH</option>
      <option value="DELETE">DELETE</option>
    </select>
  </div>

  <div class="flex flex-col gap-1.5">
    <label
      for="status-filter"
      class="text-xs font-semibold text-[var(--text-subtle)] uppercase tracking-wider"
      >Status Code</label
    >
    <input
      id="status-filter"
      type="number"
      bind:value={statusCode}
      placeholder="e.g. 200, 404, 500"
      class="h-9 px-3 w-40 rounded-lg border border-[var(--border-color)] bg-transparent text-[13px] text-[var(--text-main)] outline-none focus:border-[var(--accent)] placeholder:text-[var(--text-muted)]"
    />
  </div>

  <button
    onclick={applyFilters}
    class="h-9 px-4 rounded-lg bg-[var(--accent)] text-white text-[13px] font-semibold transition-transform hover:scale-[1.02] active:scale-[0.98]"
  >
    Apply Filters
  </button>

  {#if method || statusCode}
    <button
      onclick={clearFilters}
      class="h-9 px-4 rounded-lg bg-transparent border border-[var(--border-color)] text-[var(--text-main)] text-[13px] font-semibold hover:bg-[var(--bg-element)] transition-colors"
    >
      Clear
    </button>
  {/if}
</div>
