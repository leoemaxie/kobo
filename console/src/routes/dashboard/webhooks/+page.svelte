<script lang="ts">
  import { Plus } from '@lucide/svelte';
  import WebhookList from './WebhookList.svelte';
</script>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Page bar -->
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    padding-bottom: 20px; border-bottom: 1px solid #1e1e1e;
  ">
    <div>
      <p style="
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: #555; margin: 0 0 6px;
      ">Webhooks</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: #555;">endpoints:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: rgba(192,255,0,0.08); border: 1px solid rgba(192,255,0,0.2);
          border-radius: 4px; padding: 2px 8px; color: #C0FF00; letter-spacing: 0.05em;
        ">2 / 5</code>
      </div>
    </div>
    <button style="
      display: flex; align-items: center; gap: 6px;
      border: 1px solid #C0FF00; border-radius: 6px;
      background: #C0FF00; padding: 6px 12px;
      font-size: 13px; font-weight: 700; color: #080808; cursor: pointer;
    ">
      <Plus size={13} /> Add Endpoint
    </button>
  </div>

  <WebhookList />
</div>
