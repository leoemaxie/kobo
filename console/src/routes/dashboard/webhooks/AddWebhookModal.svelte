<script lang="ts">
  import { X, Copy } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  export let onClose: () => void;

  let secretKey = '';
  let isLoading = false;
  let selectedEvents: Record<string, boolean> = {
    'account.created': false,
    'transaction.completed': false,
    'transaction.failed': false,
    'account.closed': false
  };

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }

  $: eventsList = Object.keys(selectedEvents).filter(k => selectedEvents[k]).join(',');
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
  <div class="w-full max-w-md bg-carbon border border-iron rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-iron bg-void-black">
      <h2 class="text-sm font-semibold text-paper uppercase tracking-widest">Add Webhook Endpoint</h2>
      <button on:click={onClose} class="text-smoke hover:text-paper transition-colors">
        <X size={16} />
      </button>
    </div>

    {#if !secretKey}
      <form method="POST" action="?/addEndpoint" use:enhance={() => {
        isLoading = true;
        return async ({ result, update }) => {
          isLoading = false;
          if (result.type === 'success' && result.data?.secret) {
            secretKey = result.data.secret;
            toast.success('Webhook endpoint created');
          } else {
            toast.error(result.data?.error || 'Failed to add endpoint');
          }
          await update({ reset: false });
        };
      }} class="p-6 flex flex-col gap-5">
        
        <div>
          <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Endpoint URL</label>
          <input type="url" name="url" placeholder="https://api.yourdomain.com/webhooks" required 
                 class="w-full bg-void-black border border-iron rounded-lg px-4 py-2.5 text-sm font-mono text-paper focus:border-electric-lime focus:outline-none transition-colors" />
        </div>

        <div>
          <label class="block text-xs font-semibold text-smoke mb-2 uppercase tracking-wide">Events to Listen For</label>
          <input type="hidden" name="events" value={eventsList} />
          <div class="flex flex-col gap-2 bg-void-black border border-iron p-3 rounded-lg">
            {#each Object.keys(selectedEvents) as event}
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={selectedEvents[event]} class="rounded bg-carbon border-iron text-electric-lime focus:ring-electric-lime focus:ring-offset-void-black" />
                <span class="text-sm font-mono text-paper">{event}</span>
              </label>
            {/each}
          </div>
        </div>

        <div class="pt-2 flex justify-end gap-3">
          <button type="button" on:click={onClose} class="px-4 py-2 text-sm font-medium text-smoke hover:text-paper transition-colors">Cancel</button>
          <button type="submit" disabled={isLoading || !eventsList} class="px-5 py-2 bg-electric-lime text-void-black text-sm font-bold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50">
            {isLoading ? 'Adding...' : 'Add Endpoint'}
          </button>
        </div>
      </form>
    {:else}
      <div class="p-6 flex flex-col gap-5">
        <div class="bg-electric-lime/10 border border-electric-lime/30 rounded-lg p-4">
          <p class="text-sm font-medium text-electric-lime mb-2">Please copy this webhook secret now.</p>
          <p class="text-xs text-smoke leading-relaxed">Use this secret to verify webhook payloads. It will not be shown again.</p>
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
