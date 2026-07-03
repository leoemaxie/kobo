<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Search, Ban, CheckCircle, MoreHorizontal } from '@lucide/svelte';

  let searchQuery = $state('');
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-inter font-bold text-pure-white tracking-tight">Integrators</h1>
    <div class="relative w-64">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <Search size={16} class="text-smoke" />
      </div>
      <input
        type="text"
        bind:value={searchQuery}
        class="block w-full rounded-radius-md border border-iron bg-carbon py-2 pl-10 pr-3 text-sm text-paper placeholder-fog focus:border-steel focus:outline-none focus:ring-1 focus:ring-steel transition-colors"
        placeholder="Search integrators..."
      />
    </div>
  </div>

  <Card class="p-0 overflow-hidden">
    <table class="min-w-full divide-y divide-iron">
      <thead class="bg-carbon">
        <tr>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Integrator</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Status</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Env Access</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Joined</th>
          <th scope="col" class="relative px-6 py-3"><span class="sr-only">Actions</span></th>
        </tr>
      </thead>
      <tbody class="bg-void-black divide-y divide-iron">
        {#each [
          { name: 'Triumph Systems', email: 'dev@triumphsystems.tech', status: 'active', prodAccess: true, joined: 'Aug 12, 2026' },
          { name: 'EduPay Global', email: 'hello@edupay.example.com', status: 'active', prodAccess: false, joined: 'Sep 01, 2026' },
          { name: 'Fraudsters Inc', email: 'test@fraud.example.com', status: 'suspended', prodAccess: false, joined: 'Sep 05, 2026' }
        ] as integrator}
          <tr class="hover:bg-graphite/50 transition-colors">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex flex-col">
                <span class="text-sm font-medium text-paper">{integrator.name}</span>
                <span class="text-xs text-smoke">{integrator.email}</span>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              {#if integrator.status === 'active'}
                <span class="inline-flex items-center rounded-radius-full bg-dark-olive px-2.5 py-0.5 text-xs font-medium text-electric-lime">Active</span>
              {:else}
                <span class="inline-flex items-center rounded-radius-full bg-iron px-2.5 py-0.5 text-xs font-medium text-fog">Suspended</span>
              {/if}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center gap-2">
                <span class="h-2 w-2 rounded-full bg-smoke shadow-[0_0_8px_rgba(160,160,160,0.5)]" title="Sandbox Access"></span>
                {#if integrator.prodAccess}
                  <span class="h-2 w-2 rounded-full bg-electric-lime shadow-[0_0_8px_rgba(250,255,105,0.5)]" title="Production Access"></span>
                {/if}
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-smoke">{integrator.joined}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium flex justify-end gap-2">
              {#if !integrator.prodAccess && integrator.status === 'active'}
                <Button variant="pill" class="!px-2 !py-1 text-xs gap-1">
                  <CheckCircle size={12} /> Grant Prod
                </Button>
              {/if}
              {#if integrator.status === 'active'}
                <button class="text-smoke hover:text-red-400 transition-colors" title="Suspend Integrator">
                  <Ban size={16} />
                </button>
              {/if}
              <button class="text-smoke hover:text-electric-lime transition-colors" title="More Options">
                <MoreHorizontal size={16} />
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </Card>
</div>
