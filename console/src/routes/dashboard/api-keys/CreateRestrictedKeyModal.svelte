<script lang="ts">
  import { X, Copy } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let onClose: () => void;

  let secretKey = '';
  let isLoading = false;

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm overflow-y-auto">
  <div class="w-full max-w-md my-8 bg-carbon border border-iron rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-iron bg-void-black">
      <h2 class="text-sm font-semibold text-paper uppercase tracking-widest">Create Restricted Key</h2>
      <button on:click={onClose} class="text-smoke hover:text-paper transition-colors">
        <X size={16} />
      </button>
    </div>

    {#if !secretKey}
      <form method="POST" action="?/createKey" use:enhance={() => {
        isLoading = true;
        return async ({ result, update }) => {
          isLoading = false;
          if (result.type === 'success' && result.data?.plainSecret) {
            secretKey = result.data.plainSecret;
            toast.success('Restricted key created successfully');
          } else {
            toast.error(result.data?.error || 'Failed to create restricted key');
          }
          await update({ reset: false });
        };
      }} class="p-6 flex flex-col gap-5">
        <div>
          <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Key Label</label>
          <input type="text" name="label" placeholder="e.g. CI/CD Pipeline" required 
                 class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm text-paper focus:border-electric-lime focus:outline-none transition-colors" />
        </div>
        
        <div>
          <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Environment</label>
          <select name="environment" class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm text-paper focus:border-electric-lime focus:outline-none transition-colors">
            <option value="sandbox">Sandbox</option>
            <option value="production">Production</option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Scopes (comma-separated)</label>
          <input type="text" name="scopes" placeholder="accounts:read, transactions:write" 
                 class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm font-mono text-paper focus:border-electric-lime focus:outline-none transition-colors" />
        </div>

        <div>
          <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Allowed IPs (one per line)</label>
          <textarea name="ips" placeholder="192.168.1.1&#10;10.0.0.0/24" rows="3"
                    class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm font-mono text-paper focus:border-electric-lime focus:outline-none transition-colors"></textarea>
        </div>

        <div class="pt-2 flex justify-end gap-3">
          <button type="button" on:click={onClose} class="px-4 py-2 text-sm font-medium text-smoke hover:text-paper transition-colors">Cancel</button>
          <button type="submit" disabled={isLoading} class="px-5 py-2 bg-electric-lime text-void-black text-sm font-bold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50">
            {isLoading ? 'Creating...' : 'Create Key'}
          </button>
        </div>
      </form>
    {:else}
      <div class="p-6 flex flex-col gap-5">
        <div class="bg-electric-lime/10 border border-electric-lime/30 rounded-lg p-4">
          <p class="text-sm font-medium text-electric-lime mb-2">Please copy this secret key now.</p>
          <p class="text-xs text-smoke leading-relaxed">For security reasons, this is the only time it will be shown.</p>
        </div>
        <div class="relative">
          <input type="text" readonly value={secretKey} class="w-full bg-void-black border border-iron rounded-lg pl-4 pr-12 py-3 text-sm font-mono text-paper" />
          <button on:click={handleCopy} class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-smoke hover:text-electric-lime transition-colors">
            <Copy size={16} />
          </button>
        </div>
        <div class="pt-2 flex justify-end">
          <button on:click={onClose} class="px-5 py-2 bg-carbon border border-iron text-paper text-sm font-bold rounded-lg hover:bg-graphite transition-colors">Done</button>
        </div>
      </div>
    {/if}
  </div>
</div>
