<script lang="ts">
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let members: any[] = [];
</script>

<div
  style="background: var(--bg-element); border: 1px solid var(--border-subtle); border-radius: 8px; overflow: visible;"
>
  <div
    style="padding: 16px 20px; border-bottom: 1px solid var(--border-subtle); display: flex; align-items: center; justify-content: space-between; background: var(--bg-sidebar); border-top-left-radius: 7px; border-top-right-radius: 7px;"
  >
    <h3 style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0;">
      Workspace Members
    </h3>
  </div>
  <div style="overflow-x: auto;">
    <table style="width: 100%; min-width: 600px; border-collapse: collapse; text-align: left;">
      <thead>
        <tr style="border-bottom: 1px solid var(--border-subtle);">
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Email</th
          >
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Role</th
          >
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >MFA</th
          >
          <th
            style="padding: 12px 20px; font-size: 11px; font-weight: 600; color: var(--text-subtle); text-transform: uppercase; letter-spacing: 0.05em;"
            >Status</th
          >
          <th style="padding: 12px 20px; width: 40px;"></th>
        </tr>
      </thead>
      <tbody>
        {#each members as member}
          <tr class="team-row">
            <td
              style="padding: 16px 20px; font-family: monospace; font-size: 13px; color: var(--text-main);"
              >{member.email}</td
            >
            <td style="padding: 16px 20px;">
              <span
                style="font-family: monospace; font-size: 11px; background: var(--bg-active); border: 1px solid var(--border-color); border-radius: 4px; padding: 2px 6px; color: var(--text-muted);"
                >{member.role}</span
              >
            </td>
            <td
              style="padding: 16px 20px; font-size: 13px; color: {member.mfa
                ? 'var(--accent)'
                : 'var(--text-muted)'};">{member.mfa ? 'Enabled' : 'Disabled'}</td
            >
            <td style="padding: 16px 20px;">
              <span
                style="font-size: 13px; color: {member.status === 'Active'
                  ? 'var(--text-main)'
                  : 'var(--text-subtle)'};">{member.status}</span
              >
            </td>
            <td style="padding: 16px 20px; text-align: right; white-space: nowrap;">
              <div class="flex items-center justify-end gap-2">
                <form
                  method="POST"
                  action="?/changeRole"
                  use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Role updated');
                      else toast.error((result as any).data?.error || 'Failed to update');
                      await update();
                    };
                  }}
                >
                  <input type="hidden" name="id" value={member.id} />
                  <input
                    type="hidden"
                    name="role"
                    value={member.role === 'owner' ? 'member' : 'owner'}
                  />
                  <button
                    type="submit"
                    class="text-[12px] px-3 py-1.5 rounded-md font-medium text-main bg-element hover:bg-element-hover border border-border transition-colors"
                  >
                    Make {member.role === 'owner' ? 'Member' : 'Owner'}
                  </button>
                </form>

                <form
                  method="POST"
                  action="?/removeMember"
                  use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Member removed');
                      else toast.error((result as any).data?.error || 'Failed to remove');
                      await update();
                    };
                  }}
                >
                  <input type="hidden" name="id" value={member.id} />
                  <button
                    type="submit"
                    onclick={(e) => {
                      if (!confirm('Are you sure?')) e.preventDefault();
                    }}
                    class="text-[12px] px-3 py-1.5 rounded-md font-medium text-red-400 bg-red-400/10 hover:bg-red-400/20 border border-red-500/20 transition-colors"
                  >
                    Remove
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
  .team-row {
    border-bottom: 1px solid var(--border-subtle);
    transition: background 0.2s;
  }
  .team-row:last-child {
    border-bottom: none;
  }
  .team-row:hover {
    background: var(--bg-sidebar);
  }
  .team-row:last-child td:first-child {
    border-bottom-left-radius: 7px;
  }
  .team-row:last-child td:last-child {
    border-bottom-right-radius: 7px;
  }
</style>
