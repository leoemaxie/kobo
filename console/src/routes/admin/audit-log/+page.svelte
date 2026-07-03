<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import { Search, Filter } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const state = useConsoleState();
  const logs = $derived(state.adminAuditLogs);

  let searchQuery = $state('');
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-inter font-bold text-pure-white tracking-tight">Audit Log</h1>
      <p class="text-smoke text-sm mt-1">Immutable record of all superadmin actions.</p>
    </div>
    
    <div class="flex items-center gap-3">
      <button class="flex items-center gap-2 px-3 py-2 rounded-radius-md border border-iron bg-carbon hover:bg-graphite text-sm font-medium text-paper transition-colors">
        <Filter size={16} /> Filter
      </button>
      <div class="relative w-64">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search size={16} class="text-smoke" />
        </div>
        <input
          type="text"
          bind:value={searchQuery}
          class="block w-full rounded-radius-md border border-iron bg-carbon py-2 pl-10 pr-3 text-sm text-paper placeholder-fog focus:border-steel focus:outline-none focus:ring-1 focus:ring-steel transition-colors"
          placeholder="Search logs..."
        />
      </div>
    </div>
  </div>

  <Card class="p-0 overflow-hidden">
    <table class="min-w-full divide-y divide-iron">
      <thead class="bg-carbon">
        <tr>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Timestamp</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Actor</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Action</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Target</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Details</th>
        </tr>
      </thead>
      <tbody class="bg-void-black divide-y divide-iron text-sm text-paper">
        {#each logs as log}
          <tr class="hover:bg-graphite/50 transition-colors">
            <td class="px-6 py-4 whitespace-nowrap text-smoke font-mono text-xs">{log.time}</td>
            <td class="px-6 py-4 whitespace-nowrap font-medium">{log.actor}</td>
            <td class="px-6 py-4 whitespace-nowrap text-electric-lime">{log.action}</td>
            <td class="px-6 py-4 whitespace-nowrap">{log.target}</td>
            <td class="px-6 py-4 whitespace-nowrap text-smoke font-mono text-xs">{log.detail}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </Card>
</div>
