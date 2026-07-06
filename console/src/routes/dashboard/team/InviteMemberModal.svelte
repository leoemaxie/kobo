<script lang="ts">
  import { X } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let onClose: () => void;
  let isLoading = false;
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
  <div class="w-full max-w-md bg-carbon border border-iron rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-iron bg-void-black">
      <h2 class="text-sm font-semibold text-paper uppercase tracking-widest">Invite Team Member</h2>
      <button on:click={onClose} class="text-smoke hover:text-paper transition-colors">
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
          toast.error(result.data?.error || 'Failed to send invitation');
        }
        await update();
      };
    }} class="p-6 flex flex-col gap-5">
      
      <div>
        <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Email Address</label>
        <input type="email" name="email" placeholder="colleague@company.com" required 
               class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm text-paper focus:border-electric-lime focus:outline-none transition-colors" />
      </div>

      <div>
        <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Role</label>
        <select name="role" class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm text-paper focus:border-electric-lime focus:outline-none transition-colors">
          <option value="member">Member</option>
          <option value="owner">Owner</option>
        </select>
        <p class="mt-2 text-xs text-smoke leading-relaxed">Owners have full access to billing, settings, and team management. Members can manage API keys and webhooks.</p>
      </div>

      <div class="pt-2 flex justify-end gap-3">
        <button type="button" on:click={onClose} class="px-4 py-2 text-sm font-medium text-smoke hover:text-paper transition-colors">Cancel</button>
        <button type="submit" disabled={isLoading} class="px-5 py-2 bg-electric-lime text-void-black text-sm font-bold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50">
          {isLoading ? 'Sending...' : 'Send Invitation'}
        </button>
      </div>
    </form>
  </div>
</div>
