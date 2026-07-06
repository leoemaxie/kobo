<script lang="ts">
  import { Save } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let integrator: any;
  const state = useConsoleState();
</script>

<form method="POST" action="?/updateWorkspace" use:enhance={() => {
  return async ({ result, update }) => {
    if (result.type === 'failure') {
      toast.error(result.data?.error as string || 'Update failed.');
    } else if (result.type === 'error') {
      toast.error('An unexpected server error occurred.');
    } else {
      toast.success('Workspace updated successfully.');
    }
    await update();
  };
}} style="background: var(--bg-element); border: 1px solid var(--border-subtle); border-radius: 8px; overflow: hidden;">
  <div style="padding: 16px 20px; border-bottom: 1px solid var(--border-subtle); background: var(--bg-sidebar);">
    <h3 style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0;">General Information</h3>
    <p style="font-size: 12px; color: var(--text-subtle); margin: 4px 0 0;">Manage your workspace details.</p>
  </div>
  <div style="padding: 20px; display: grid; gap: 16px;">
    <div>
      <label for="workspaceName" style="display: block; font-size: 12px; font-weight: 600; color: var(--text-muted); margin-bottom: 6px;">Workspace Name</label>
      <input id="workspaceName" type="text" name="name" value={integrator?.name || 'Kobo Inc.'} class="settings-input" required />
    </div>
    <div>
      <label for="supportEmail" style="display: block; font-size: 12px; font-weight: 600; color: var(--text-muted); margin-bottom: 6px;">Support Email</label>
      <input id="supportEmail" type="email" value={state.user?.email || ''} class="settings-input" disabled />
      <p style="font-size: 11px; color: var(--text-subtle); margin: 6px 0 0;">Linked to the owner's authentication email.</p>
    </div>
    <div style="margin-top: 10px;">
      <button type="submit" style="
        display: inline-flex; align-items: center; gap: 6px;
        border: 1px solid var(--accent); border-radius: 6px;
        background: var(--accent); padding: 6px 12px;
        font-size: 13px; font-weight: 700; color: var(--accent-text); cursor: pointer;
      ">
        <Save size={13} /> Save Changes
      </button>
    </div>
  </div>
</form>

<style>
  .settings-input {
    width: 100%;
    max-width: 400px;
    background: var(--bg-element);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 8px 12px;
    font-size: 13px;
    color: var(--text-main);
    outline: none;
    transition: border-color 0.2s;
  }
  .settings-input:focus {
    border-color: var(--accent);
  }
</style>
