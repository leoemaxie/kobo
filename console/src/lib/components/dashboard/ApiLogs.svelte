<script lang="ts">
  import { ChevronRight, ExternalLink } from '@lucide/svelte';

  const logs = [
    { method: 'POST', path: '/v1/accounts',              status: 201, ms: 87,  id: 'req_9Kz2', time: 'just now' },
    { method: 'GET',  path: '/v1/identities/id_892nf8',  status: 200, ms: 43,  id: 'req_8mXp', time: '2m ago'  },
    { method: 'POST', path: '/v1/transactions',           status: 400, ms: 120, id: 'req_7nQr', time: '15m ago' },
    { method: 'GET',  path: '/v1/accounts?limit=10',      status: 200, ms: 31,  id: 'req_6wKs', time: '1h ago'  },
    { method: 'DELETE',path:'/v1/hooks/wh_92nd',          status: 204, ms: 55,  id: 'req_5jVt', time: '3h ago'  },
  ];

  const methodColor: Record<string, string> = {
    GET:    '#888',
    POST:   '#C0FF00',
    DELETE: '#f87171',
    PATCH:  '#60a5fa',
  };

  function statusColor(s: number) {
    if (s < 300) return '#C0FF00';
    if (s < 400) return '#60a5fa';
    return '#f87171';
  }
</script>

<div>
  <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;">
    <p style="font-size: 12px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #555;">
      Recent Requests
    </p>
    <a href="/logs" style="
      font-size: 12px; font-weight: 600; color: #555; text-decoration: none;
      display: flex; align-items: center; gap: 4px;
    "
      onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.color = '#C0FF00'}
      onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.color = '#555'}
    >
      View all <ExternalLink size={11} />
    </a>
  </div>

  <div style="border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
    <!-- Table header -->
    <div style="
      display: grid; grid-template-columns: 60px 1fr 60px 70px 80px 32px;
      padding: 8px 16px; background: #0d0d0d; border-bottom: 1px solid #1e1e1e;
    ">
      {#each ['METHOD','PATH','STATUS','TIME','REQ ID',''] as col}
        <span style="font-size: 11px; font-weight: 700; letter-spacing: 0.1em; color: #444; text-transform: uppercase;">
          {col}
        </span>
      {/each}
    </div>

    {#each logs as log}
      <div style="
        display: grid; grid-template-columns: 60px 1fr 60px 70px 80px 32px;
        padding: 10px 16px; border-bottom: 1px solid #111;
        align-items: center; cursor: pointer; transition: background 0.1s;
      "
        onmouseenter={(e) => (e.currentTarget as HTMLDivElement).style.background = '#0f0f0f'}
        onmouseleave={(e) => (e.currentTarget as HTMLDivElement).style.background = 'transparent'}
      >
        <span style="
          font-family: monospace; font-size: 12px; font-weight: 700;
          color: {methodColor[log.method] ?? '#888'};
        ">{log.method}</span>

        <span style="
          font-family: 'JetBrains Mono', monospace; font-size: 13px;
          color: #C8C8C8; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        ">{log.path}</span>

        <span style="
          font-family: monospace; font-size: 13px; font-weight: 600;
          color: {statusColor(log.status)};
        ">{log.status}</span>

        <span style="font-family: monospace; font-size: 13px; color: #555;">{log.ms}ms</span>

        <span style="
          font-family: monospace; font-size: 12px; color: #444;
        ">{log.id}</span>

        <ChevronRight size={13} color="#333" />
      </div>
    {/each}
  </div>
</div>
