<script lang="ts">
  import { Copy, RefreshCw, Plus, Trash2 } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import RollKeyModal from './RollKeyModal.svelte';
  import { createEventDispatcher } from 'svelte';

  export let keys: any[] = [];
  const dispatch = createEventDispatcher();

  let rollingKeyId: string | null = null;

  const cols = ['NAME', 'KEY ID', 'SECRET KEY', 'LAST USED', 'CREATED', ''];
</script>

<div>
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    margin-bottom: 10px;
  ">
    <p style="
      font-size: 10px; font-weight: 700; text-transform: uppercase;
      letter-spacing: 0.1em; color: var(--text-subtle); margin: 0;
    ">Standard Keys</p>
    <button on:click={() => dispatch('create')} style="
      display: flex; align-items: center; gap: 5px;
      border: 1px solid #2a2a2a; border-radius: 6px;
      background: var(--bg-sidebar); padding: 5px 10px;
      font-size: 11px; font-weight: 600; color: var(--text-muted); cursor: pointer;
    ">
      <Plus size={12} /> Create secret key
    </button>
  </div>

  <div style="border: 1px solid var(--border-subtle); border-radius: 8px; overflow: hidden;">
    <!-- Header -->
    <div style="
      display: grid; grid-template-columns: 1fr 1.4fr 1.6fr 90px 100px 72px;
      padding: 9px 16px; background: var(--bg-sidebar); border-bottom: 1px solid var(--border-subtle);
    ">
      {#each cols as col}
        <span style="font-size: 11px; font-weight: 700; letter-spacing: 0.1em; color: var(--text-muted); text-transform: uppercase;">
          {col}
        </span>
      {/each}
    </div>

    {#each keys as k}
      <div style="
        display: grid; grid-template-columns: 1fr 1.4fr 1.6fr 90px 100px 72px;
        padding: 11px 16px; align-items: center; border-bottom: 1px solid var(--bg-sidebar);
      ">
        <span style="font-size: 14px; font-weight: 500; color: #C8C8C8;">{k.name}</span>

        <code style="font-family: monospace; font-size: 13px; color: var(--text-muted);">{k.id}</code>

        <code style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">
          ••••••••••••••••••
        </code>

        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">{k.lastUsed}</span>
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">{k.created}</span>

        <div style="display: flex; align-items: center; gap: 10px; justify-content: flex-end;">
          <button on:click={() => { navigator.clipboard.writeText(k.id); toast.success('Key ID copied'); }} 
            style="background: none; border: none; cursor: pointer; color: var(--text-subtle); padding: 0; display: flex;"
            title="Copy key ID">
            <Copy size={13} />
          </button>
          <button on:click={() => rollingKeyId = k.id}
            style="background: none; border: none; cursor: pointer; color: var(--text-subtle); padding: 0; display: flex;"
            title="Roll key">
            <RefreshCw size={13} />
          </button>
          <form method="POST" action="?/revokeKey" use:enhance={() => {
            return async ({ result, update }) => {
              if (result.type === 'success') {
                toast.success('Key revoked');
              } else {
                toast.error((result as any).data?.error || 'Failed to revoke key');
              }
              await update();
            };
          }} style="display: inline-block;">
            <input type="hidden" name="keyId" value={k.id} />
            <button type="submit" style="background: none; border: none; cursor: pointer; color: var(--text-subtle); padding: 0; display: flex;"
              title="Revoke key" on:click={(e) => { if(!confirm('Are you sure you want to revoke this key?')) e.preventDefault(); }}>
              <Trash2 size={13} />
            </button>
          </form>
        </div>
      </div>
    {/each}
  </div>
</div>

{#if rollingKeyId}
  <RollKeyModal keyId={rollingKeyId} onClose={() => rollingKeyId = null} />
{/if}
