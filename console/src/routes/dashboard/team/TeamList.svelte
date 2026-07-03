<script lang="ts">
  import { MoreVertical } from '@lucide/svelte';
  const members = [
    { email: 'admin@kobo.inc', role: 'Owner', status: 'Active', mfa: true },
    { email: 'developer@kobo.inc', role: 'Admin', status: 'Active', mfa: true },
    { email: 'finance@kobo.inc', role: 'Viewer', status: 'Pending', mfa: false }
  ];
</script>

<div style="background: #0a0a0a; border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
  <div style="padding: 16px 20px; border-bottom: 1px solid #1e1e1e; display: flex; align-items: center; justify-content: space-between; background: #111;">
    <h3 style="font-size: 14px; font-weight: 600; color: #F8F8F8; margin: 0;">Workspace Members</h3>
  </div>
  <table style="width: 100%; border-collapse: collapse; text-align: left;">
    <thead>
      <tr style="border-bottom: 1px solid #1e1e1e;">
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Email</th>
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Role</th>
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">MFA</th>
        <th style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: #666; text-transform: uppercase; letter-spacing: 0.05em;">Status</th>
        <th style="padding: 12px 20px; width: 40px;"></th>
      </tr>
    </thead>
    <tbody>
      {#each members as member}
        <tr class="team-row">
          <td style="padding: 16px 20px; font-family: monospace; font-size: 13px; color: #F8F8F8;">{member.email}</td>
          <td style="padding: 16px 20px;">
            <span style="font-family: monospace; font-size: 11px; background: #1a1a1a; border: 1px solid #333; border-radius: 4px; padding: 2px 6px; color: #aaa;">{member.role}</span>
          </td>
          <td style="padding: 16px 20px; font-size: 13px; color: {member.mfa ? '#C0FF00' : '#888'};">{member.mfa ? 'Enabled' : 'Disabled'}</td>
          <td style="padding: 16px 20px;">
            <span style="font-size: 13px; color: {member.status === 'Active' ? '#F8F8F8' : '#666'};">{member.status}</span>
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
  .team-row {
    border-bottom: 1px solid #1e1e1e;
    transition: background 0.2s;
  }
  .team-row:hover {
    background: #111;
  }
</style>
