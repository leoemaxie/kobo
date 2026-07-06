<script lang="ts">
  import WorkspaceInfoCard from './WorkspaceInfoCard.svelte';
  import DangerZoneCard from './DangerZoneCard.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import CodeBadge from '$lib/components/ui/CodeBadge.svelte';

  export let data;
</script>

<svelte:head>
  <title>Settings — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="Workspace Settings">
    {#snippet meta()}
      <span class="font-inconsolata text-[13px] text-subtle">id:</span>
      <CodeBadge>{data.integrator?.id}</CodeBadge>
    {/snippet}
  </PageHeader>

  <div class="grid gap-6">
    <WorkspaceInfoCard integrator={data.integrator} />
    <DangerZoneCard />
  </div>
</div>
