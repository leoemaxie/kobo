<script lang="ts">
  import { X } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let onClose: () => void;
  let isLoading = false;
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
  <div class="w-full max-w-md bg-element border border-border rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-border bg-sidebar">
      <h2 class="text-sm font-semibold text-main uppercase tracking-widest">Invite Team Member</h2>
      <button onclick={onClose} class="text-muted hover:text-main transition-colors">
        <X size={16} />
      </button>
    </div>

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
      
      <div>
        <label for="email" class="block text-xs font-semibold text-muted mb-2 uppercase tracking-wide">Email Address</label>
        <input id="email" type="email" name="email" placeholder="colleague@company.com" required 
               class="w-full bg-sidebar border border-border rounded-lg px-4 py-2.5 text-sm text-main focus:border-border-focus focus:outline-none transition-colors" />
      </div>

      <div>
        <label for="role" class="block text-xs font-semibold text-muted mb-2 uppercase tracking-wide">Role</label>
        <select id="role" name="role" class="w-full bg-sidebar border border-border rounded-lg px-4 py-2.5 text-sm text-main focus:border-border-focus focus:outline-none transition-colors">
          <option value="member">Member</option>
          <option value="owner">Owner</option>
        </select>
        <p class="mt-2 text-xs text-muted leading-relaxed">Owners have full access to billing, settings, and team management. Members can manage API keys and webhooks.</p>
      </div>

      <div class="pt-2 flex justify-end gap-3">
        <button type="button" onclick={onClose} class="px-4 py-2 text-sm font-medium text-muted hover:text-main transition-colors">Cancel</button>
        <button type="submit" disabled={isLoading} class="px-5 py-2 bg-primary text-primary-text text-sm font-bold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50">
          {isLoading ? 'Sending...' : 'Send Invitation'}
        </button>
      </div>
    </form>
  </div>
</div>
