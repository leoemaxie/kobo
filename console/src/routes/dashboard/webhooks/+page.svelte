<script lang="ts">
  import { Plus } from '@lucide/svelte';
  import WebhookList from './WebhookList.svelte';
  import AddWebhookModal from './AddWebhookModal.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import CodeBadge from '$lib/components/ui/CodeBadge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  let { data } = $props();

  const consoleState = useConsoleState();
  let currentEnv = $derived(consoleState.currentEnvironment);
  let filteredEndpoints = $derived(data.endpoints.filter((e: any) => e.environment === currentEnv || !e.environment));

  let showAddModal = $state(false);
</script>

<svelte:head>
  <title>Webhooks — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="Webhooks">
    {#snippet meta()}
      <span class="font-inconsolata text-[13px] text-subtle">environment:</span>
      <CodeBadge>{currentEnv}</CodeBadge>
      <span class="text-[11px] text-muted">·</span>
      <span class="font-inconsolata text-[13px] text-subtle">endpoints:</span>
      <CodeBadge variant="neutral">{filteredEndpoints.length} / 5</CodeBadge>
    {/snippet}
    {#snippet actions()}
      <Button variant="primary" size="md" onclick={() => showAddModal = true}>
        <Plus size={13} /> Add Endpoint
      </Button>
    {/snippet}
  </PageHeader>

  <WebhookList endpoints={filteredEndpoints} />
</div>

{#if showAddModal}
  <AddWebhookModal onClose={() => showAddModal = false} />
{/if}
