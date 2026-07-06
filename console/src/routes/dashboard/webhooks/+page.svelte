<script lang="ts">
  import { Plus } from '@lucide/svelte';
  import WebhookList from './WebhookList.svelte';
  import AddWebhookModal from './AddWebhookModal.svelte';

  export let data;

  let showAddModal = false;
</script>

<svelte:head>
  <title>Webhooks — Kobo Console</title>
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
      ">Webhooks</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">endpoints:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: var(--accent-transparent); border: 1px solid var(--accent-border);
          border-radius: 4px; padding: 2px 8px; color: var(--accent); letter-spacing: 0.05em;
        ">{data.endpoints.length} / 5</code>
      </div>
    </div>
    <button on:click={() => showAddModal = true} style="
      display: flex; align-items: center; gap: 6px;
      border: 1px solid var(--accent); border-radius: 7px;
      background: var(--accent); padding: 6px 12px;
      font-size: 13px; font-weight: 700; color: var(--accent-text); cursor: pointer;
      transition: all 0.15s;
    ">
      <Plus size={13} /> Add Endpoint
    </button>
  </div>

  <WebhookList endpoints={data.endpoints} />
</div>
