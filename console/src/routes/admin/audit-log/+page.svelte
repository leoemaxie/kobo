<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import { Search, Filter } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const consoleState = useConsoleState();
  const logs = $derived(consoleState.adminAuditLogs);

  let searchQuery = $state('');
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-inter font-bold text-main tracking-tight">Audit Log</h1>
      <p class="text-muted text-sm mt-1">Immutable record of all superadmin actions.</p>
    </div>
    
    <div class="flex items-center gap-3">
      <button class="flex items-center gap-2 px-3 py-2 rounded-radius-md border border-border bg-sidebar hover:bg-element-hover text-sm font-medium text-main transition-colors">
        <Filter size={16} /> Filter
      </button>
      <div class="relative w-64">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search size={16} class="text-muted" />
        </div>
        <input
          type="text"
          bind:value={searchQuery}
          class="block w-full rounded-radius-md border border-border bg-sidebar py-2 pl-10 pr-3 text-sm text-main placeholder-fog focus:border-border-focus focus:outline-none focus:ring-1 focus:ring-border-focus transition-colors"
          placeholder="Search logs..."
        />
      </div>
    </div>
  </div>

  <Card class="p-0 overflow-hidden">
    <table class="min-w-full divide-y divide-border">
      <thead class="bg-sidebar">
        <tr>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Timestamp</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Actor</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Action</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Target</th>
          <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-muted uppercase">Details</th>
        </tr>
      </thead>
      <tbody class="bg-element divide-y divide-border text-sm text-main">
        {#each logs as log}
          <tr class="hover:bg-element-hover transition-colors">
            <td class="px-6 py-4 whitespace-nowrap text-muted font-mono text-xs">{log.time}</td>
            <td class="px-6 py-4 whitespace-nowrap font-medium">{log.actor}</td>
            <td class="px-6 py-4 whitespace-nowrap text-primary">{log.action}</td>
            <td class="px-6 py-4 whitespace-nowrap">{log.target}</td>
            <td class="px-6 py-4 whitespace-nowrap text-muted font-mono text-xs">{log.detail}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </Card>
</div>
