<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const consoleState = useConsoleState();
  const logs = $derived(consoleState.adminAuditLogs);

</script>

<svelte:head>
  <title>Audit Log | Kobo Console</title>
</svelte:head>


<div class="space-y-6">

  <Card class="p-0 overflow-hidden">
    <div class="overflow-x-auto">
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
    </div>
  </Card>
</div>
