<script lang="ts">
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  export let onClose: () => void;
  let isLoading = false;
</script>

<Modal title="Invite Team Member" {onClose}>
  <form method="POST" action="?/inviteMember" use:enhance={() => {
    isLoading = true;
    return async ({ result, update }) => {
      isLoading = false;
      if (result.type === 'success') {
        toast.success('Invitation sent');
        onClose();
      } else {
        toast.error((result as any).data?.error || 'Failed to send invitation');
      }
      await update();
    };
  }} class="p-6 flex flex-col gap-5">
    
    <Input id="email" label="Email Address" type="email" name="email" placeholder="colleague@company.com" required />

    <div class="space-y-1.5">
      <label for="role" class="block text-xs font-semibold text-muted uppercase tracking-widest">Role</label>
      <select id="role" name="role" class="w-full bg-sidebar border border-border rounded-lg px-4 py-2.5 text-sm text-main focus:border-border-focus focus:outline-none transition-colors">
        <option value="member">Member</option>
        <option value="owner">Owner</option>
      </select>
      <p class="text-xs text-muted leading-relaxed">Owners have full access to billing, settings, and team management. Members can manage API keys and webhooks.</p>
    </div>

    <div class="pt-2 flex justify-end gap-3">
      <Button type="button" variant="ghost" size="md" onclick={onClose}>Cancel</Button>
      <Button type="submit" variant="primary" size="md" disabled={isLoading}>
        {isLoading ? 'Sending...' : 'Send Invitation'}
      </Button>
    </div>
  </form>
</Modal>
