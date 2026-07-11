<script lang="ts">
  import LogFilters from '$lib/components/dashboard/LogFilters.svelte';
  import LogTable from '$lib/components/dashboard/LogTable.svelte';
  import PaginationControls from '$lib/components/dashboard/PaginationControls.svelte';
  import { ArrowLeft } from '@lucide/svelte';

  let { data } = $props<{
    data: {
      paginatedLogs: any[];
      meta: { total: number; page: number; limit: number; totalPages: number };
      filters: { method: string; statusCode: string };
    };
  }>();

  let logs = $derived(data.paginatedLogs);
  let meta = $derived(data.meta);
  let filters = $derived(data.filters);
</script>

<svelte:head>
  <title>API Logs | Kobo Console</title>
</svelte:head>

<div class="max-w-6xl mx-auto">
  <div class="mb-8">
    <a
      href="/dashboard"
      class="inline-flex items-center gap-1.5 text-[13px] font-semibold text-[var(--text-subtle)] hover:text-[var(--text-main)] transition-colors mb-4"
    >
      <ArrowLeft size={14} />
      Back to Dashboard
    </a>
    <h1 class="text-2xl font-bold text-[var(--text-main)] tracking-tight">API Logs</h1>
    <p class="text-[14px] text-[var(--text-subtle)] mt-1">
      Browse and filter through your recent API requests.
    </p>
  </div>

  <LogFilters />

  <LogTable {logs} />

  <PaginationControls {meta} />
</div>
