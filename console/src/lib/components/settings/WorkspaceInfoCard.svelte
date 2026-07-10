<script lang="ts">
  import { Save } from '@lucide/svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import CardSection from '$lib/components/ui/CardSection.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  export let integrator: any;
  const state = useConsoleState();
</script>

<form
  method="POST"
  action="?/updateWorkspace"
  use:enhance={() => {
    return async ({ result, update }) => {
      if (result.type === 'failure') {
        toast.error((result.data?.error as string) || 'Update failed.');
      } else if (result.type === 'error') {
        toast.error('An unexpected server error occurred.');
      } else {
        toast.success('Workspace updated successfully.');
      }
      await update();
    };
  }}
>
  <CardSection title="General Information" subtitle="Manage your workspace details.">
    <div class="grid gap-4">
      <Input
        id="workspaceName"
        label="Workspace Name"
        type="text"
        name="name"
        value={integrator?.name || 'Kobo Inc.'}
        variant="settings"
        required
      />
      <Input
        id="supportEmail"
        label="Support Email"
        type="email"
        value={state.user?.email || ''}
        variant="settings"
        disabled
        hint="Linked to the owner's authentication email."
      />
      <div class="mt-2">
        <Button type="submit" variant="primary" size="md">
          <Save size={13} /> Save Changes
        </Button>
      </div>
    </div>
  </CardSection>
</form>
