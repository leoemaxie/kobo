<script lang="ts">
  import { MoreVertical, Circle, Eye, EyeOff } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let endpoints: any[] = [];

  let secretRevealed: Record<string, boolean> = {};
  let openDropdownId: string | null = null;
</script>

<div
  style="background: var(--bg-element); border: 1px solid var(--border-subtle); border-radius: 8px; overflow: hidden;"
>
  <div
    style="padding: 16px 20px; border-bottom: 1px solid var(--border-subtle); display: flex; align-items: center; justify-content: space-between; background: var(--bg-sidebar);"
  >
    <h3 style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0;">
      Configured Endpoints
    </h3>
  </div>
  <div style="overflow-x: auto;">
    <table style="width: 100%; min-width: 800px; border-collapse: collapse; text-align: left;">
      <thead>
        <tr style="border-bottom: 1px solid var(--border-subtle);">
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Endpoint URL</th
          >
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Events</th
          >
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Secret</th
          >
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Status</th
          >
          <th style="padding: 12px 20px; width: 40px;"></th>
        </tr>
      </thead>
      <tbody>
        {#each endpoints as ep}
          <tr class="webhook-row">
            <td
              style="padding: 16px 20px; font-family: monospace; font-size: 13px; color: var(--text-main);"
              >{ep.url}</td
            >
            <td style="padding: 16px 20px;">
              <div style="display: flex; gap: 6px; flex-wrap: wrap;">
                {#each ep.events as ev}
                  <span
                    style="font-family: monospace; font-size: 11px; background: var(--bg-active); border: 1px solid var(--border-color); border-radius: 4px; padding: 2px 6px; color: var(--text-muted);"
                    >{ev}</span
                  >
                {/each}
              </div>
            </td>
            <td style="padding: 16px 20px;">
              <div style="display: flex; align-items: center; gap: 8px;">
                <code
                  style="font-family: monospace; font-size: 13px; color: {secretRevealed[ep.id]
                    ? 'var(--text-main)'
                    : 'var(--text-muted)'};"
                >
                  {secretRevealed[ep.id] ? ep.secret : '••••••••••••••••••'}
                </code>
                <button
                  onclick={() => (secretRevealed[ep.id] = !secretRevealed[ep.id])}
                  style="background: none; border: none; cursor: pointer; color: var(--text-subtle); padding: 0; display: flex;"
                >
                  {#if secretRevealed[ep.id]}<EyeOff size={13} />{:else}<Eye size={13} />{/if}
                </button>
              </div>
            </td>
            <td style="padding: 16px 20px;">
              <div style="display: flex; align-items: center; gap: 6px;">
                <Circle
                  size={8}
                  color={ep.status === 'active' ? 'var(--accent)' : 'var(--text-muted)'}
                  fill={ep.status === 'active' ? 'var(--accent)' : 'var(--text-muted)'}
                />
                <span
                  style="font-size: 13px; color: {ep.status === 'active'
                    ? 'var(--text-main)'
                    : 'var(--text-subtle)'}; text-transform: capitalize;">{ep.status}</span
                >
              </div>
            </td>
            <td style="padding: 16px 20px;">
              <div style="display: flex; gap: 8px; justify-content: flex-end; align-items: center;">
                <form
                  method="POST"
                  action="?/toggleEndpoint"
                  use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Status updated');
                      else toast.error((result as any).data?.error || 'Failed to update');
                      await update();
                    };
                  }}
                >
                  <input type="hidden" name="id" value={ep.id} />
                  <input type="hidden" name="currentStatus" value={ep.status} />
                  <button
                    type="submit"
                    class="hover:bg-[var(--bg-active)] transition-colors"
                    style="font-size: 12px; font-weight: 500; color: var(--text-main); background: transparent; border: 1px solid var(--border-color); border-radius: 6px; padding: 5px 10px; cursor: pointer;"
                  >
                    {ep.status === 'active' ? 'Disable' : 'Enable'}
                  </button>
                </form>

                <form
                  method="POST"
                  action="?/deleteEndpoint"
                  use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Webhook deleted');
                      else toast.error((result as any).data?.error || 'Failed to delete');
                      await update();
                    };
                  }}
                >
                  <input type="hidden" name="id" value={ep.id} />
                  <button
                    type="submit"
                    onclick={(e) => {
                      if (!confirm('Are you sure you want to delete this webhook?'))
                        e.preventDefault();
                    }}
                    class="hover:bg-red-500/10 transition-colors"
                    style="font-size: 12px; font-weight: 500; color: #ef4444; background: transparent; border: 1px solid rgba(239, 68, 68, 0.2); border-radius: 6px; padding: 5px 10px; cursor: pointer;"
                  >
                    Delete
                  </button>
                </form>
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<style>
  .webhook-row {
    border-bottom: 1px solid var(--border-subtle);
    transition: background 0.2s;
  }
  .webhook-row:hover {
    background: var(--bg-sidebar);
  }
</style>
