<script lang="ts">
  import { Plus } from '@lucide/svelte';
  import TeamList from './TeamList.svelte';
  import InviteMemberModal from '../../../lib/components/team/InviteMemberModal.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import CodeBadge from '$lib/components/ui/CodeBadge.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  export let data;

  let showInviteModal = false;
</script>

<svelte:head>
  <title>Teams — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="Team Management">
    {#snippet meta()}
      <span class="font-inconsolata text-[13px] text-subtle">members:</span>
      <CodeBadge>{data.members.length} / 10</CodeBadge>
    {/snippet}
    {#snippet actions()}
      <Button variant="primary" size="md" onclick={() => showInviteModal = true}>
        <Plus size={13} /> Invite Member
      </Button>
    {/snippet}
  </PageHeader>

  <TeamList members={data.members} />
</div>

{#if showInviteModal}
  <InviteMemberModal onClose={() => showInviteModal = false} />
{/if}
