<script lang="ts">
  import { ShieldAlert, Plus, Lock } from '@lucide/svelte';
  import { createEventDispatcher } from 'svelte';
  import SectionLabel from '$lib/components/ui/SectionLabel.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  
  export let keys: any[] = [];
  const dispatch = createEventDispatcher();
</script>

<div>
  <SectionLabel>Restricted Keys</SectionLabel>

  <div class="border border-dashed border-border rounded-lg px-5 py-7 flex items-center justify-between gap-6">
    <div class="flex items-start gap-3.5">
      <ShieldAlert size={16} class="text-subtle mt-0.5 shrink-0" />
      <div>
        <p class="text-xs font-semibold text-muted mb-1">
          Restricted Keys — not configured
        </p>
        <p class="text-[11px] text-subtle leading-relaxed">
          Scope keys to specific IP ranges and permission sets for service-to-service use.
        </p>
        <div class="flex items-center gap-1.5 mt-2">
          <Lock size={11} class="text-muted" />
          <span class="font-mono text-[10px] text-muted">scopes: accounts:read, transactions:write</span>
        </div>
      </div>
    </div>

    <Button variant="neutral" size="sm" class="shrink-0 whitespace-nowrap" on:click={() => dispatch('create')}>
      <Plus size={12} /> Create restricted key
    </Button>
  </div>

  {#if keys.length > 0}
    <div class="mt-2.5 flex flex-col gap-2">
      {#each keys as k}
        <div class="border border-border rounded-[6px] px-4 py-3 flex items-center justify-between">
          <div>
            <p class="text-[13px] font-medium text-main mb-1">{k.name}</p>
            <code class="font-mono text-[11px] text-muted">{k.id}</code>
          </div>
          <span class="text-[11px] text-subtle">Created {k.created}</span>
        </div>
      {/each}
    </div>
  {/if}
</div>
