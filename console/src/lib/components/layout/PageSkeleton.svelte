<script lang="ts">
  import Skeleton from '$lib/components/ui/Skeleton.svelte';
</script>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Header Skeleton -->
  <div
    style="display: flex; justify-content: space-between; border-bottom: 1px solid var(--border-color); padding-bottom: 32px;"
  >
    <div style="display: flex; flex-direction: column; gap: 12px;">
      <Skeleton style="height: 24px; width: 140px;" />
      <Skeleton style="height: 16px; width: 250px;" />
    </div>
    <div style="display: flex; gap: 8px;">
      <Skeleton style="height: 32px; width: 80px;" />
      <Skeleton style="height: 32px; width: 120px;" />
    </div>
  </div>

  <!-- Metric row -->
  <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px;">
    <Skeleton style="height: 120px; width: 100%; border-radius: 10px;" />
    <Skeleton style="height: 120px; width: 100%; border-radius: 10px;" />
    <Skeleton style="height: 120px; width: 100%; border-radius: 10px;" />
    <Skeleton style="height: 120px; width: 100%; border-radius: 10px;" />
  </div>

  <!-- Content Block -->
  <Skeleton style="height: 400px; width: 100%; border-radius: 8px;" />
</div>
