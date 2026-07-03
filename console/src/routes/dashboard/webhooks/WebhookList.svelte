<script lang="ts">
  import { MoreVertical, Circle } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const endpoints = $derived(state.webhooks);
</script>

<div style="background: #0a0a0a; border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
  <div style="padding: 16px 20px; border-bottom: 1px solid #1e1e1e; display: flex; align-items: center; justify-content: space-between; background: #111;">
    <h3 style="font-size: 14px; font-weight: 600; color: #F8F8F8; margin: 0;">Configured Endpoints</h3>
  </div>
  <table style="width: 100%; border-collapse: collapse; text-align: left;">
    <thead>
      <tr style="border-bottom: 1px solid #1e1e1e;">
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Endpoint URL</th>
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Events</th>
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Secret</th>
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Status</th>
        <th style="padding: 12px 20px; width: 40px;"></th>
      </tr>
    </thead>
    <tbody>
      {#each endpoints as ep}
        <tr class="webhook-row">
          <td style="padding: 16px 20px; font-family: monospace; font-size: 13px; color: #F8F8F8;">{ep.url}</td>
          <td style="padding: 16px 20px;">
            <div style="display: flex; gap: 6px; flex-wrap: wrap;">
              {#each ep.events as ev}
                <span style="font-family: monospace; font-size: 11px; background: #1a1a1a; border: 1px solid #333; border-radius: 4px; padding: 2px 6px; color: #aaa;">{ev}</span>
              {/each}
            </div>
          </td>
          <td style="padding: 16px 20px; font-family: monospace; font-size: 13px; color: #888;">{ep.secret}</td>
          <td style="padding: 16px 20px;">
            <div style="display: flex; align-items: center; gap: 6px;">
              <Circle size={8} color={ep.status === 'active' ? '#C0FF00' : '#444'} fill={ep.status === 'active' ? '#C0FF00' : '#444'} />
              <span style="font-size: 13px; color: {ep.status === 'active' ? '#F8F8F8' : '#666'}; text-transform: capitalize;">{ep.status}</span>
            </div>
          </td>
          <td style="padding: 16px 20px; text-align: right;">
            <button style="background: transparent; border: none; color: #666; cursor: pointer; padding: 4px;"><MoreVertical size={16} /></button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .webhook-row {
    border-bottom: 1px solid #1e1e1e;
    transition: background 0.2s;
  }
  .webhook-row:hover {
    background: #111;
  }
</style>
