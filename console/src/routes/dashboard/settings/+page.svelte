<script lang="ts">
  import { Save, AlertTriangle } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import { enhance } from '$app/forms';

  let { data } = $props();
  const state = useConsoleState();
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
      ">Workspace Settings</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: #555;">id:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: rgba(192,255,0,0.08); border: 1px solid rgba(192,255,0,0.2);
          border-radius: 4px; padding: 2px 8px; color: #C0FF00; letter-spacing: 0.05em;
        ">{state.user?.integratorId}</code>
      </div>
    </div>
    <!-- We will use the form submit button in the body instead -->
  </div>

  <div style="display: grid; gap: 24px;">
    <!-- General Settings -->
    <form method="POST" action="?/updateWorkspace" use:enhance style="background: #0a0a0a; border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
      <div style="padding: 16px 20px; border-bottom: 1px solid #1e1e1e; background: #111;">
        <h3 style="font-size: 14px; font-weight: 600; color: #F8F8F8; margin: 0;">General Information</h3>
        <p style="font-size: 12px; color: #666; margin: 4px 0 0;">Manage your workspace details.</p>
      </div>
      <div style="padding: 20px; display: grid; gap: 16px;">
        <div>
          <label style="display: block; font-size: 12px; font-weight: 600; color: #888; margin-bottom: 6px;">Workspace Name</label>
          <input type="text" name="name" value={data.integrator?.name || 'Kobo Inc.'} class="settings-input" required />
        </div>
        <div>
          <label style="display: block; font-size: 12px; font-weight: 600; color: #888; margin-bottom: 6px;">Support Email</label>
          <input type="email" value={state.user?.email || ''} class="settings-input" disabled />
          <p style="font-size: 11px; color: #555; margin: 6px 0 0;">Linked to the owner's authentication email.</p>
        </div>
        <div style="margin-top: 10px;">
          <button type="submit" style="
            display: inline-flex; align-items: center; gap: 6px;
            border: 1px solid #C0FF00; border-radius: 6px;
            background: #C0FF00; padding: 6px 12px;
            font-size: 13px; font-weight: 700; color: #080808; cursor: pointer;
          ">
            <Save size={13} /> Save Changes
          </button>
        </div>
      </div>
    </form>

    <!-- Danger Zone -->
    <div style="background: #0a0a0a; border: 1px solid #3a1c1c; border-radius: 8px; overflow: hidden;">
      <div style="padding: 16px 20px; border-bottom: 1px solid #3a1c1c; background: #1c0f0f;">
        <h3 style="font-size: 14px; font-weight: 600; color: #ff4a4a; margin: 0;">Danger Zone</h3>
      </div>
      <div style="padding: 20px; display: flex; align-items: center; justify-content: space-between;">
        <div>
          <h4 style="font-size: 13px; font-weight: 600; color: #F8F8F8; margin: 0 0 4px;">Delete Workspace</h4>
          <p style="font-size: 12px; color: #888; margin: 0;">Permanently remove this workspace and all associated data.</p>
        </div>
        <form method="POST" action="?/deleteWorkspace">
          <button type="submit" class="delete-btn" onclick="return confirm('Are you sure? This action cannot be undone and will delete all API Keys, Webhooks, and Billing records.')">
            <AlertTriangle size={13} /> Delete Workspace
          </button>
        </form>
      </div>
    </div>
  </div>
</div>

<style>
  .settings-input {
    width: 100%;
    max-width: 400px;
    background: #111;
    border: 1px solid #2a2a2a;
    border-radius: 6px;
    padding: 8px 12px;
    font-size: 13px;
    color: #F8F8F8;
    outline: none;
    transition: border-color 0.2s;
  }
  .settings-input:focus {
    border-color: #C0FF00;
  }

  .delete-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: #1a0a0a;
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
    color: #000;
  }
</style>
