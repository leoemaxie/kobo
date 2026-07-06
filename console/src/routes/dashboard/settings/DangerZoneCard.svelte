<script lang="ts">
  import { AlertTriangle } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
</script>

<div class="danger-zone-container">
  <div class="danger-zone-header">
    <h3 style="font-size: 14px; font-weight: 600; color: #ff4a4a; margin: 0;">Danger Zone</h3>
  </div>
  <div style="padding: 20px; display: flex; align-items: center; justify-content: space-between;">
    <div>
      <h4 style="font-size: 13px; font-weight: 600; color: var(--text-main); margin: 0 0 4px;">Delete Workspace</h4>
      <p style="font-size: 12px; color: var(--text-muted); margin: 0;">Permanently remove this workspace and all associated data.</p>
    </div>
    <form method="POST" action="?/deleteWorkspace" use:enhance={({ cancel }) => {
      if (!confirm('Are you sure? This action cannot be undone and will delete all API Keys, Webhooks, and Billing records.')) cancel();
      return async ({ result, update }) => {
        if (result.type === 'failure') {
          toast.error(result.data?.error as string || 'Deletion failed.');
        } else if (result.type === 'error') {
          toast.error('An unexpected server error occurred.');
        } else {
          toast.success('Workspace deleted.');
        }
        await update();
      };
    }}>
      <button type="submit" class="delete-btn">
        <AlertTriangle size={13} /> Delete Workspace
      </button>
    </form>
  </div>
</div>

<style>
  .delete-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    border: 1px solid #ff4a4a;
    border-radius: 6px;
    padding: 6px 12px;
    font-size: 13px;
    font-weight: 600;
    color: #ff4a4a;
    cursor: pointer;
    transition: all 0.2s;
  }
  .delete-btn:hover {
    background: #ff4a4a;
    color: #fff;
  }
  :global(:root[data-theme='dark']) .delete-btn:hover {
    color: #000;
  }

  .danger-zone-container {
    background: var(--bg-element);
    border: 1px solid #ff4a4a40;
    border-radius: 8px;
    overflow: hidden;
  }
  .danger-zone-header {
    padding: 16px 20px;
    border-bottom: 1px solid #ff4a4a40;
    background: #ff4a4a10;
  }
</style>
