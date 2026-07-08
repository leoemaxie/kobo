<script lang="ts">
  import { AlertTriangle, LogOut } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import CardSection from '$lib/components/ui/CardSection.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { page } from '$app/stores';
</script>

<CardSection title="Danger Zone" danger>
  <div class="flex flex-col gap-6">
    <!-- Leave Workspace -->
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h4 class="text-[13px] font-semibold text-main mb-1">Leave Workspace</h4>
        <p class="text-xs text-muted">Remove your account from this workspace. You will lose access immediately.</p>
      </div>
      <form
        method="POST"
        action="?/leaveWorkspace"
        use:enhance={({ cancel }) => {
          if (!confirm('Are you sure you want to leave this workspace? You will need an invitation to rejoin.')) cancel();
          return async ({ result, update }) => {
            if (result.type === 'failure') {
              toast.error(result.data?.error as string || 'Failed to leave workspace.');
            } else if (result.type === 'error') {
              toast.error('An unexpected server error occurred.');
            } else {
              toast.success('You have left the workspace.');
            }
            await update();
          };
        }}
      >
        <Button type="submit" variant="secondary" size="md">
          <LogOut size={13} /> Leave Workspace
        </Button>
      </form>
    </div>

    <!-- Delete Workspace (Only visible to owners) -->
    {#if $page.data.user?.role === 'owner'}
      <div class="w-full h-[1px] bg-border-subtle my-2"></div>
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
    {/if}
  </div>
</CardSection>
