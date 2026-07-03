<script lang="ts">
  import { BookOpen, ArrowUpRight, Terminal } from '@lucide/svelte';
  import {PUBLIC_DOCS_URL} from '$env/static/public';
</script>

<div class="flex items-start justify-between gap-6 pb-8 border-b border-iron">
  <div class="space-y-3">
    <div class="flex items-center gap-3">
      <h1 style="
        font-size: 15px; font-weight: 600; color: #888;
        text-transform: uppercase; letter-spacing: 0.08em;
      ">Overview</h1>
      <span style="
        display: inline-flex; align-items: center; gap: 5px;
        background: rgba(192,255,0,0.08); border: 1px solid rgba(192,255,0,0.2);
        border-radius: 4px; padding: 2px 7px;
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: #C0FF00;
      ">Sandbox</span>
    </div>
    <div style="display: flex; align-items: baseline; gap: 8px;">
      <span style="
        font-family: 'JetBrains Mono', monospace; font-size: 13px; color: #555;
      ">project_id:</span>
      <code style="
        font-family: 'JetBrains Mono', monospace; font-size: 13px; color: #C8C8C8;
        background: #111; border: 1px solid #2a2a2a; border-radius: 4px; padding: 2px 7px;
      ">proj_sandbox_a4b9c12f</code>
    </div>
  </div>

  <div style="display: flex; align-items: center; gap: 8px;">
    <a
      href={PUBLIC_DOCS_URL}
      style="
        display: flex; align-items: center; gap: 6px;
        border: 1px solid #2a2a2a; border-radius: 7px;
        background: #111; padding: 6px 12px;
        font-size: 13px; font-weight: 500; color: #888; text-decoration: none;
        transition: all 0.15s;
      "
      onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.color = '#F8F8F8'}
      onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.color = '#888'}
    >
      <BookOpen size={13} /> Docs
    </a>
    <a
      href="/dashboard/api-keys"
      style="
        display: flex; align-items: center; gap: 6px;
        border: 1px solid #C0FF00; border-radius: 7px;
        background: #C0FF00; padding: 6px 12px;
        font-size: 13px; font-weight: 700; color: #080808; text-decoration: none;
        transition: all 0.15s;
      "
    >
      <Terminal size={13} /> Get API Key
      <ArrowUpRight size={12} />
    </a>
  </div>
</div>
