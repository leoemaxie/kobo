<script lang="ts">
  import { Plus } from '@lucide/svelte';
  import StandardKeysTable from '$lib/components/api-keys/StandardKeysTable.svelte';
  import RestrictedKeysSection from '$lib/components/api-keys/RestrictedKeysSection.svelte';
  import CreateKeyModal from '$lib/components/api-keys/CreateKeyModal.svelte';
  import CreateRestrictedKeyModal from '$lib/components/api-keys/CreateRestrictedKeyModal.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import CodeBadge from '$lib/components/ui/CodeBadge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  let { data } = $props();

  const consoleState = useConsoleState();
  let currentEnv = $derived(consoleState.currentEnvironment);

  let showCreateModal = $state(false);
  let showCreateRestrictedModal = $state(false);
</script>

<svelte:head>
  <title>API Keys — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="API Keys">
    {#snippet meta()}
      <span class="font-inconsolata text-[13px] text-subtle">environment:</span>
      <CodeBadge>{currentEnv}</CodeBadge>
      <span class="text-[11px] text-muted">·</span>
      <span class="text-[13px] text-subtle">Toggle in header to switch to production</span>
    {/snippet}
    {#snippet actions()}
      <Button variant="primary" size="md" onclick={() => (showCreateModal = true)}>
        <Plus size={13} /> Create Key
      </Button>
    {/snippet}
  </PageHeader>

  <StandardKeysTable
    keys={data.keys.filter(
      (k) => k.status === 'active' && !k.id.includes('restricted') && k.environment === currentEnv,
    )}
    on:create={() => (showCreateModal = true)}
  />
  <RestrictedKeysSection
    keys={data.keys.filter(
      (k) => k.status === 'active' && k.id.includes('restricted') && k.environment === currentEnv,
    )}
    on:create={() => (showCreateRestrictedModal = true)}
  />
</div>

{#if showCreateModal}
  <CreateKeyModal onClose={() => (showCreateModal = false)} />
{/if}

{#if showCreateRestrictedModal}
  <CreateRestrictedKeyModal onClose={() => (showCreateRestrictedModal = false)} />
{/if}
