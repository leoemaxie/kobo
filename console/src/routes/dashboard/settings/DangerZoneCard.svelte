<script lang="ts">
  import { AlertTriangle } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import CardSection from '$lib/components/ui/CardSection.svelte';
  import Button from '$lib/components/ui/Button.svelte';
</script>

<CardSection title="Danger Zone" danger>
  <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
    <div>
      <h4 class="text-[13px] font-semibold text-main mb-1">Delete Workspace</h4>
      <p class="text-xs text-muted">Permanently remove this workspace and all associated data.</p>
    </div>
    <form
      method="POST"
      action="?/deleteWorkspace"
      use:enhance={({ cancel }) => {
        if (!confirm('Are you sure? This action cannot be undone and will delete all API Keys, Webhooks, and Billing records.')) cancel();
        return async ({ result, update }) => {
          if (result.type === 'failure') {
            toast.error(result.data?.error as string || 'Deletion failed.');
          } else if (result.type === 'error') {
            toast.error('An unexpected server error occurred.');
          } else {
            toast.success('Workspace deleted.');
          }
          await update();
        };
      }}
    >
      <Button type="submit" variant="danger" size="md">
        <AlertTriangle size={13} /> Delete Workspace
      </Button>
    </form>
  </div>
</CardSection>
