<script lang="ts">
  import { ChevronRight, ExternalLink } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const consoleState = useConsoleState();
  const logs = $derived(consoleState.logs);

  const methodColor: Record<string, string> = {
    GET:    'var(--text-muted)',
    POST:   'var(--accent)',
    DELETE: 'var(--error-color)',
    PATCH:  '#60a5fa',
  };

  function statusColor(s: number) {
    if (s < 300) return 'var(--accent)';
    if (s < 400) return '#60a5fa';
    return 'var(--error-color)';
  }
</script>

<div>
  <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;">
    <p style="font-size: 12px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle);">
      Recent Requests
    </p>
    <a href="/logs" style="
      font-size: 12px; font-weight: 600; color: var(--text-subtle); text-decoration: none;
      display: flex; align-items: center; gap: 4px; transition: color 0.2s;
    "
      onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.color = 'var(--accent)'}
      onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.color = 'var(--text-subtle)'}
    >
      View all <ExternalLink size={11} />
    </a>
  </div>

  <div style="border: 1px solid var(--border-color); border-radius: 8px; overflow: hidden;">
    <!-- Table header -->
    <div style="
      display: grid; grid-template-columns: 60px 1fr 60px 70px 80px 32px;
      padding: 8px 16px; background: var(--bg-active); border-bottom: 1px solid var(--border-color);
    ">
      {#each ['METHOD','PATH','STATUS','TIME','REQ ID',''] as col}
        <span style="font-size: 11px; font-weight: 700; letter-spacing: 0.1em; color: var(--text-muted); text-transform: uppercase;">
          {col}
        </span>
      {/each}
    </div>

    {#each logs as log}
      <div role="row" tabindex="0" style="
        display: grid; grid-template-columns: 60px 1fr 60px 70px 80px 32px;
        padding: 10px 16px; border-bottom: 1px solid var(--border-color);
        align-items: center; cursor: pointer; transition: background 0.1s;
      "
        onmouseenter={(e) => (e.currentTarget as HTMLDivElement).style.background = 'var(--bg-element)'}
        onmouseleave={(e) => (e.currentTarget as HTMLDivElement).style.background = 'transparent'}
      >
        <span style="
          font-family: monospace; font-size: 12px; font-weight: 700;
          color: {methodColor[log.method] ?? 'var(--text-muted)'};
        ">{log.method}</span>

        <span style="
          font-family: 'JetBrains Mono', monospace; font-size: 13px;
          color: var(--text-main); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        ">{log.path}</span>

        <span style="
          font-family: monospace; font-size: 13px; font-weight: 600;
          color: {statusColor(log.status)};
        ">{log.status}</span>

        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">{log.ms}ms</span>

        <span style="
          font-family: monospace; font-size: 12px; color: var(--text-muted);
        ">{log.id}</span>

        <ChevronRight size={13} color="var(--text-subtle)" />
      </div>
    {/each}
  </div>
</div>
