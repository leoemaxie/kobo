<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Search, Ban, CheckCircle, Play } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let data;
  const integrators = data.integrators;

  let searchQuery = '';
</script>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
    <h1 class="text-2xl font-inter font-bold text-main tracking-tight">Integrators</h1>
    <div class="relative w-full md:w-64">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <Search size={16} class="text-muted" />
      </div>
      <input
        type="text"
        bind:value={searchQuery}
        class="block w-full rounded-radius-md border border-border bg-sidebar py-2 pl-10 pr-3 text-sm text-main placeholder-fog focus:border-border-focus focus:outline-none focus:ring-1 focus:ring-border-focus transition-colors"
        placeholder="Search integrators..."
      />
    </div>
  </div>

  <Card class="p-0 overflow-hidden">
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-border">
        <thead class="bg-sidebar">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Integrator</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Status</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Env Access</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Joined</th>
            <th scope="col" class="relative px-6 py-3"><span class="sr-only">Actions</span></th>
          </tr>
        </thead>
        <tbody class="bg-element divide-y divide-border">
          {#each integrators as integrator}
            <tr class="hover:bg-element-hover transition-colors">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex flex-col">
                  <span class="text-sm font-medium text-main">{integrator.name}</span>
                  <span class="text-xs text-muted">{integrator.email}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                {#if integrator.status === 'active'}
                  <span class="inline-flex items-center rounded-radius-full bg-primary-transparent border border-primary-border px-2.5 py-0.5 text-xs font-medium text-primary">Active</span>
                {:else}
                  <span class="inline-flex items-center rounded-radius-full bg-sidebar border border-border px-2.5 py-0.5 text-xs font-medium text-subtle">Suspended</span>
                {/if}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <span class="h-2 w-2 rounded-full bg-muted shadow-sm" title="Sandbox Access"></span>
                  {#if integrator.prodAccess}
                    <span class="h-2 w-2 rounded-full bg-primary shadow-sm" title="Production Access"></span>
                  {/if}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-muted">{integrator.joined}</td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium flex justify-end gap-2">
                {#if !integrator.prodAccess && integrator.status === 'active'}
                  <form method="POST" action="?/grantProduction" use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Granted production access');
                      else toast.error((result as any).data?.error || 'Failed to grant access');
                      await update();
                    };
                  }}>
                    <input type="hidden" name="id" value={integrator.id} />
                    <Button type="submit" variant="pill" class="!px-2 !py-1 text-xs gap-1">
                      <CheckCircle size={12} /> Grant Prod
                    </Button>
                  </form>
                {/if}
                
                {#if integrator.status === 'active'}
                  <form method="POST" action="?/suspendIntegrator" use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Integrator suspended');
                      else toast.error((result as any).data?.error || 'Failed to suspend');
                      await update();
                    };
                  }}>
                    <input type="hidden" name="id" value={integrator.id} />
                    <button type="submit" onclick={(e) => { if(!confirm('Suspend integrator?')) e.preventDefault(); }} class="text-muted hover:text-error transition-colors mt-1.5" title="Suspend Integrator">
                      <Ban size={16} />
                    </button>
                  </form>
                {:else}
                  <form method="POST" action="?/reinstateIntegrator" use:enhance={() => {
                    return async ({ result, update }) => {
                      if (result.type === 'success') toast.success('Integrator reinstated');
                      else toast.error((result as any).data?.error || 'Failed to reinstate');
                      await update();
                    };
                  }}>
                    <input type="hidden" name="id" value={integrator.id} />
                    <button type="submit" class="text-muted hover:text-primary transition-colors mt-1.5" title="Reinstate Integrator">
                      <Play size={16} />
                    </button>
                  </form>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </Card>
</div>
