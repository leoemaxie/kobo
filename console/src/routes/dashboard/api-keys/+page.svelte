<script lang="ts">
  import { Plus } from '@lucide/svelte';
  import StandardKeysTable from './StandardKeysTable.svelte';
  import RestrictedKeysSection from './RestrictedKeysSection.svelte';
  import CreateKeyModal from './CreateKeyModal.svelte';
  import CreateRestrictedKeyModal from './CreateRestrictedKeyModal.svelte';
  
  import { useConsoleState } from '$lib/state/console.svelte';

  let { data } = $props();

  const consoleState = useConsoleState();
  let currentEnv = $derived(consoleState.currentEnvironment);

  let showCreateModal = $state(false);
  let showCreateRestrictedModal = $state(false);
</script>

<svelte:head>
  <title>API Keys — Kobo Console</title>
</svelte:head>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Page bar -->
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    padding-bottom: 20px; border-bottom: 1px solid var(--border-subtle);
  ">
    <div>
      <p style="
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 6px;
      ">API Keys</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="
          font-family: monospace; font-size: 13px; color: var(--text-subtle);
        ">environment:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: var(--accent-transparent); border: 1px solid var(--accent-border);
          border-radius: 4px; padding: 2px 8px; color: var(--accent); letter-spacing: 0.05em;
        ">{currentEnv}</code>
        <span style="font-size: 11px; color: var(--text-muted);">·</span>
        <span style="font-size: 13px; color: var(--text-subtle);">Toggle in header to switch to production</span>
      </div>
    </div>
    <button onclick={() => showCreateModal = true} style="
      display: flex; align-items: center; gap: 6px;
      border: 1px solid var(--accent); border-radius: 7px;
      background: var(--accent); padding: 6px 12px;
      font-size: 13px; font-weight: 700; color: var(--accent-text); cursor: pointer;
      transition: all 0.15s;
    ">
      <Plus size={13} /> Create Key
    </button>
  </div>

  <StandardKeysTable keys={data.keys.filter(k => k.status === 'active' && !k.id.includes('restricted') && k.environment === currentEnv)} on:create={() => showCreateModal = true} />
  <RestrictedKeysSection keys={data.keys.filter(k => k.status === 'active' && k.id.includes('restricted') && k.environment === currentEnv)} on:create={() => showCreateRestrictedModal = true} />
</div>

{#if showCreateModal}
  <CreateKeyModal onClose={() => showCreateModal = false} />
{/if}

{#if showCreateRestrictedModal}
  <CreateRestrictedKeyModal onClose={() => showCreateRestrictedModal = false} />
{/if}
