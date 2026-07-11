<script lang="ts">
  import PageHeader from '$lib/components/common/PageHeader.svelte';
  import AdminList from '$lib/components/admin/super/AdminList.svelte';

  let { data, form } = $props();
  let admins = $derived(data.admins);
</script>

<svelte:head>
  <title>Superadmin | Triumph Academy</title>
</svelte:head>

{#snippet superBadge()}
  <div class="bg-carbon border border-electric-lime text-electric-lime text-xs font-bold px-3 py-1.5 rounded-full tracking-widest uppercase shadow-[0_0_10px_rgba(204,255,0,0.2)]">
    Superadmin
  </div>
{/snippet}

<div class="space-y-10 w-full">
  <PageHeader
    title="Superadmin Console"
    subtitle="Manage administrators, grant access, and configure scopes."
    rightSection={superBadge}
  />

  <AdminList {admins} {form} />
</div>
