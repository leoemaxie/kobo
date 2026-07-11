<script lang="ts">
  let { logs } = $props<{ logs: any[] }>();

  const methodColor: Record<string, string> = {
    GET: 'var(--text-muted)',
    POST: 'var(--accent)',
    DELETE: 'var(--error-color)',
    PATCH: '#60a5fa',
    PUT: '#f59e0b',
  };

  function statusColor(s: number) {
    if (s < 300) return 'var(--accent)';
    if (s < 400) return '#60a5fa';
    return 'var(--error-color)';
  }
</script>

<div class="border border-[var(--border-color)] rounded-xl overflow-x-auto bg-[var(--bg-sidebar)]">
  <div class="min-w-[800px]">
    <!-- Table header -->
    <div
      class="grid grid-cols-[80px_1fr_80px_100px_120px] px-5 py-3 bg-[var(--bg-active)] border-b border-[var(--border-color)]"
    >
      {#each ['METHOD', 'PATH', 'STATUS', 'TIME', 'REQ ID'] as col}
        <span class="text-[11px] font-bold tracking-[0.1em] text-[var(--text-muted)] uppercase"
          >{col}</span
        >
      {/each}
    </div>

    <!-- Table body -->
    {#if logs.length === 0}
      <div class="p-12 text-center">
        <p class="text-[14px] text-[var(--text-subtle)]">
          No requests found matching the criteria.
        </p>
      </div>
    {:else}
      {#each logs as log}
        <div
          class="grid grid-cols-[80px_1fr_80px_100px_120px] px-5 py-3 border-b border-[var(--border-color)] items-center transition-colors hover:bg-[var(--bg-element)] last:border-0"
        >
          <span
            class="font-mono text-[12px] font-bold"
            style="color: {methodColor[log.method] ?? 'var(--text-muted)'}">{log.method}</span
          >
          <span
            class="font-mono text-[13px] text-[var(--text-main)] whitespace-nowrap overflow-hidden text-ellipsis pr-4"
            >{log.path}</span
          >
          <span class="font-mono text-[13px] font-semibold" style="color: {statusColor(log.status)}"
            >{log.status}</span
          >
          <div class="flex flex-col gap-0.5">
            <span class="font-mono text-[13px] text-[var(--text-main)]">{log.time}</span>
            <span class="font-mono text-[11px] text-[var(--text-subtle)]">{log.ms}ms</span>
          </div>
          <span class="font-mono text-[12px] text-[var(--text-muted)]">{log.id}</span>
        </div>
      {/each}
    {/if}
  </div>
</div>
