<script lang="ts">
  import { ShieldAlert, Plus, Lock } from '@lucide/svelte';
  import { createEventDispatcher } from 'svelte';
  
  export let keys: any[] = [];
  const dispatch = createEventDispatcher();
</script>

<div>
  <p style="
    font-size: 10px; font-weight: 700; text-transform: uppercase;
    letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 10px;
  ">Restricted Keys</p>

  <div style="
    border: 1px dashed #2a2a2a; border-radius: 8px;
    padding: 28px 20px;
    display: flex; align-items: center; justify-content: space-between; gap: 24px;
  ">
    <div style="display: flex; align-items: flex-start; gap: 14px;">
      <ShieldAlert size={16} color="var(--text-subtle)" style="margin-top: 2px; flex-shrink: 0;" />
      <div>
        <p style="font-size: 12px; font-weight: 600; color: var(--text-muted); margin: 0 0 4px;">
          Restricted Keys — not configured
        </p>
        <p style="font-size: 11px; color: var(--text-subtle); margin: 0; line-height: 1.6;">
          Scope keys to specific IP ranges and permission sets for service-to-service use.
        </p>
        <div style="display: flex; align-items: center; gap: 6px; margin-top: 8px;">
          <Lock size={11} color="var(--text-muted)" />
          <span style="font-family: monospace; font-size: 10px; color: var(--text-muted);">
            scopes: accounts:read, transactions:write
          </span>
        </div>
      </div>
    </div>

    <button on:click={() => dispatch('create')} style="
      display: flex; align-items: center; gap: 6px; flex-shrink: 0;
      border: 1px solid #2a2a2a; border-radius: 6px;
      background: var(--bg-sidebar); padding: 6px 12px;
      font-size: 11px; font-weight: 600; color: var(--text-muted); cursor: pointer;
      white-space: nowrap;
    ">
      <Plus size={12} /> Create restricted key
    </button>
  </div>

  {#if keys.length > 0}
    <div style="margin-top: 10px; display: flex; flex-direction: column; gap: 8px;">
      {#each keys as k}
        <div style="border: 1px solid #2a2a2a; border-radius: 6px; padding: 12px 16px; display: flex; align-items: center; justify-content: space-between;">
          <div>
            <p style="font-size: 13px; font-weight: 500; color: #C8C8C8; margin: 0 0 4px;">{k.name}</p>
            <code style="font-family: monospace; font-size: 11px; color: var(--text-muted);">{k.id}</code>
          </div>
          <div style="display: flex; gap: 6px; align-items: center;">
             <span style="font-size: 11px; color: var(--text-subtle);">Created {k.created}</span>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
