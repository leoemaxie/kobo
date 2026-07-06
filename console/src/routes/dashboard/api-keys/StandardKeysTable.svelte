<script lang="ts">
  import { Copy, RefreshCw, Plus, Trash2 } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import RollKeyModal from './RollKeyModal.svelte';
  import SectionLabel from '$lib/components/ui/SectionLabel.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { createEventDispatcher } from 'svelte';

  export let keys: any[] = [];
  const dispatch = createEventDispatcher();

  let rollingKeyId: string | null = null;

  const cols = ['NAME', 'KEY ID', 'SECRET KEY', 'LAST USED', 'CREATED', ''];
</script>

<div>
  <div class="flex items-center justify-between mb-2.5">
    <SectionLabel class="mb-0">Standard Keys</SectionLabel>
    <Button variant="neutral" size="sm" on:click={() => dispatch('create')}>
      <Plus size={12} /> Create secret key
    </Button>
  </div>

  <div class="border border-border-subtle rounded-lg overflow-x-auto">
    <div class="min-w-[800px]">
      <!-- Header -->
    <div class="grid gap-0 bg-sidebar border-b border-border-subtle px-4 py-2.5"
         style="grid-template-columns: 1fr 1.4fr 1.6fr 90px 100px 72px;">
      {#each cols as col}
        <span class="text-[11px] font-bold tracking-widest text-muted uppercase">{col}</span>
      {/each}
    </div>

    {#each keys as k}
      <div class="grid items-center px-4 py-3 border-b border-background last:border-b-0 hover:bg-sidebar transition-colors"
           style="grid-template-columns: 1fr 1.4fr 1.6fr 90px 100px 72px;">
        <span class="text-sm font-medium text-main">{k.name}</span>
        <code class="font-mono text-[13px] text-muted">{k.id}</code>
        <code class="font-mono text-[13px] text-subtle">••••••••••••••••••</code>
        <span class="font-mono text-[13px] text-subtle">{k.lastUsed}</span>
        <span class="font-mono text-[13px] text-subtle">{k.created}</span>

        <div class="flex items-center gap-2.5 justify-end">
          <button
            on:click={() => { navigator.clipboard.writeText(k.id); toast.success('Key ID copied'); }}
            class="text-subtle hover:text-muted transition-colors"
            title="Copy key ID"
          >
            <Copy size={13} />
          </button>
          <button
            on:click={() => rollingKeyId = k.id}
            class="text-subtle hover:text-muted transition-colors"
            title="Roll key"
          >
            <RefreshCw size={13} />
          </button>
          <form method="POST" action="?/revokeKey" use:enhance={() => {
            return async ({ result, update }) => {
              if (result.type === 'success') toast.success('Key revoked');
              else toast.error((result as any).data?.error || 'Failed to revoke key');
              await update();
            };
          }} class="inline-flex">
            <input type="hidden" name="keyId" value={k.id} />
            <button
              type="submit"
              class="text-subtle hover:text-red-400 transition-colors"
              title="Revoke key"
              on:click={(e) => { if(!confirm('Are you sure you want to revoke this key?')) e.preventDefault(); }}
            >
              <Trash2 size={13} />
            </button>
          </form>
        </div>
      </div>
    {/each}
    </div>
  </div>
</div>

{#if rollingKeyId}
  <RollKeyModal keyId={rollingKeyId} onClose={() => rollingKeyId = null} />
{/if}
