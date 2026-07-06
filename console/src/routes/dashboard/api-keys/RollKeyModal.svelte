<script lang="ts">
  import { X, Copy, AlertTriangle } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let keyId: string;
  export let onClose: () => void;

  let secretKey = '';
  let isLoading = false;

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
  <div class="w-full max-w-md bg-carbon border border-iron rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-iron bg-void-black">
      <h2 class="text-sm font-semibold text-paper uppercase tracking-widest">Roll API Key</h2>
      <button on:click={onClose} class="text-smoke hover:text-paper transition-colors">
        <X size={16} />
      </button>
    </div>

    {#if !secretKey}
      <form method="POST" action="?/rollKey" use:enhance={() => {
        isLoading = true;
        return async ({ result, update }) => {
          isLoading = false;
          if (result.type === 'success' && result.data?.plainSecret) {
            secretKey = result.data.plainSecret;
            toast.success('Key rolled successfully');
          } else {
            toast.error(result.data?.error || 'Failed to roll key');
          }
          await update({ reset: false });
        };
      }} class="p-6 flex flex-col gap-5">
        <input type="hidden" name="keyId" value={keyId} />
        
        <div class="bg-red-500/10 border border-red-500/30 rounded-lg p-4 flex gap-3">
          <AlertTriangle class="text-red-400 shrink-0 mt-0.5" size={18} />
          <div>
            <p class="text-sm font-medium text-red-400 mb-1">This will invalidate the current key.</p>
            <p class="text-xs text-smoke leading-relaxed">Any applications currently using <code class="font-mono text-[10px] bg-void-black px-1 py-0.5 rounded text-paper">{keyId}</code> will stop working immediately. A new key will be generated with the same label and permissions.</p>
          </div>
        </div>

        <div class="pt-2 flex justify-end gap-3">
          <button type="button" on:click={onClose} class="px-4 py-2 text-sm font-medium text-smoke hover:text-paper transition-colors">Cancel</button>
          <button type="submit" disabled={isLoading} class="px-5 py-2 bg-red-500 text-white text-sm font-bold rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50">
            {isLoading ? 'Rolling...' : 'Roll Key'}
          </button>
        </div>
      </form>
    {:else}
      <div class="p-6 flex flex-col gap-5">
        <div class="bg-electric-lime/10 border border-electric-lime/30 rounded-lg p-4">
          <p class="text-sm font-medium text-electric-lime mb-2">Please copy your NEW secret key.</p>
          <p class="text-xs text-smoke leading-relaxed">Update your applications with this new secret immediately.</p>
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
