<script lang="ts">
  import DashboardHeader from '$lib/components/dashboard/DashboardHeader.svelte';
  import MetricCards from '$lib/components/dashboard/MetricCards.svelte';
  import ApiLogs from '$lib/components/dashboard/ApiLogs.svelte';
  import QuickActions from '$lib/components/dashboard/QuickActions.svelte';
</script>

<div class="space-y-6">
  <!-- Page header + project metadata -->
  <DashboardHeader />

  <!-- Metric row -->
  <MetricCards />

  <!-- Logs + sidebar panel -->
  <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
    <div class="xl:col-span-2">
      <ApiLogs />
    </div>
    <div>
      <QuickActions />
    </div>
  </div>
</div>
