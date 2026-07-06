<script lang="ts">
  import { BookOpen, ArrowUpRight, Terminal } from '@lucide/svelte';
  import {PUBLIC_DOCS_URL} from '$env/static/public';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
</script>

<div class="flex items-start justify-between gap-6 pb-8 border-b" style="border-color: var(--border-color);">
  <div class="space-y-3">
    <div class="flex items-center gap-3">
      <h1 style="
        font-size: 15px; font-weight: 600; color: var(--text-muted);
        text-transform: uppercase; letter-spacing: 0.08em;
      ">Overview</h1>
      <span style="
        display: inline-flex; align-items: center; gap: 5px;
        background: var(--accent-transparent); border: 1px solid var(--accent-border);
        border-radius: 4px; padding: 2px 7px;
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: var(--accent);
      ">{state.currentEnvironment}</span>
    </div>
    <div style="display: flex; align-items: baseline; gap: 8px;">
      <span style="
        font-family: 'JetBrains Mono', monospace; font-size: 13px; color: var(--text-subtle);
      ">project_id:</span>
      <code style="
        font-family: 'JetBrains Mono', monospace; font-size: 13px; color: var(--text-main);
        background: var(--bg-element); border: 1px solid var(--border-color); border-radius: 4px; padding: 2px 7px;
      ">{state.user?.integratorId?.split('-')[0] || 'proj_sandbox_a4b9c12f'}</code>
    </div>
  </div>

  <div style="display: flex; align-items: center; gap: 8px;">
    <a
      href={PUBLIC_DOCS_URL}
      style="
        display: flex; align-items: center; gap: 6px;
        border: 1px solid var(--border-color); border-radius: 7px;
        background: var(--bg-element); padding: 6px 12px;
        font-size: 13px; font-weight: 500; color: var(--text-muted); text-decoration: none;
        transition: all 0.15s;
      "
      onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.color = 'var(--text-main)'}
      onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.color = 'var(--text-muted)'}
    >
      <BookOpen size={13} /> Docs
    </a>
    <a
      href="/dashboard/api-keys"
      style="
        display: flex; align-items: center; gap: 6px;
        border: 1px solid var(--accent); border-radius: 7px;
        background: var(--accent); padding: 6px 12px;
        font-size: 13px; font-weight: 700; color: var(--accent-text); text-decoration: none;
        transition: all 0.15s;
      "
    >
      <Terminal size={13} /> Get API Key
      <ArrowUpRight size={12} />
    </a>
  </div>
</div>
