<script lang="ts">
  import WorkspaceInfoCard from './WorkspaceInfoCard.svelte';
  import DangerZoneCard from './DangerZoneCard.svelte';

  export let data;
</script>

<svelte:head>
  <title>Settings — Kobo Console</title>
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
      ">Workspace Settings</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">id:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: var(--accent-transparent); border: 1px solid var(--accent-border);
          border-radius: 4px; padding: 2px 8px; color: var(--accent); letter-spacing: 0.05em;
        ">{data.integrator?.id}</code>
      </div>
    </div>
    <!-- We will use the form submit button in the body instead -->
  </div>

  <div style="display: grid; gap: 24px;">
    <WorkspaceInfoCard integrator={data.integrator} />
    <DangerZoneCard />
  </div>
</div>
