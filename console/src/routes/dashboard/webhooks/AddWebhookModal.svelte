<script lang="ts">
  import { X, Copy } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  let { onClose }: { onClose: () => void } = $props();

  const consoleState = useConsoleState();
  let currentEnv = $derived(consoleState.currentEnvironment);

  let secretKey = $state('');
  let isLoading = $state(false);
  let selectedEvents: Record<string, boolean> = $state({
    'account.created': false,
    'transaction.completed': false,
    'transaction.failed': false,
    'account.closed': false
  });

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }

  let eventsList = $derived(Object.keys(selectedEvents).filter(k => selectedEvents[k]).join(','));
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
  <div class="w-full max-w-md bg-element border border-border rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-border bg-sidebar">
      <h2 class="text-sm font-semibold text-main uppercase tracking-widest">Add Webhook Endpoint</h2>
      <button onclick={onClose} class="text-muted hover:text-main transition-colors">
        <X size={16} />
      </button>
    </div>

    {#if !secretKey}
      <form method="POST" action="?/addEndpoint" use:enhance={() => {
        isLoading = true;
        return async ({ result, update }) => {
          isLoading = false;
          if (result.type === 'success' && (result as any).data?.secret) {
            secretKey = (result as any).data.secret as string;
            toast.success('Webhook endpoint created');
          } else {
            toast.error((result as any).data?.error || 'Failed to add endpoint');
          }
          await update({ reset: false });
        };
      }} class="p-6 flex flex-col gap-5">
        
        <input type="hidden" name="environment" value={currentEnv} />
        
        <div>
          <label for="url" class="block text-xs font-semibold text-muted mb-2 uppercase tracking-wide">Endpoint URL</label>
          <input id="url" type="url" name="url" placeholder="https://api.yourdomain.com/webhooks" required 
                 class="w-full bg-sidebar border border-border rounded-lg px-4 py-2.5 text-sm font-mono text-main focus:border-border-focus focus:outline-none transition-colors" />
        </div>

        <div>
          <label for="events" class="block text-xs font-semibold text-muted mb-2 uppercase tracking-wide">Events to Listen For</label>
          <input id="events" type="hidden" name="events" value={eventsList} />
          <div class="flex flex-col gap-2 bg-sidebar border border-border p-3 rounded-lg">
            {#each Object.keys(selectedEvents) as event}
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={selectedEvents[event]} class="rounded bg-element border-border text-primary focus:ring-primary focus:ring-offset-sidebar" />
                <span class="text-sm font-mono text-main">{event}</span>
              </label>
            {/each}
          </div>
        </div>

        <div class="pt-2 flex justify-end gap-3">
          <button type="button" onclick={onClose} class="px-4 py-2 text-sm font-medium text-muted hover:text-main transition-colors">Cancel</button>
          <button type="submit" disabled={isLoading || !eventsList} class="px-5 py-2 bg-primary text-primary-text text-sm font-bold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50">
            {isLoading ? 'Adding...' : 'Add Endpoint'}
          </button>
        </div>
      </form>
    {:else}
      <div class="p-6 flex flex-col gap-5">
        <div class="bg-primary-transparent border border-primary-border rounded-lg p-4">
          <p class="text-sm font-medium text-primary mb-2">Please copy this webhook secret now.</p>
          <p class="text-xs text-muted leading-relaxed">Use this secret to verify webhook payloads. It will not be shown again.</p>
        </div>
        <div class="relative">
          <input type="text" readonly value={secretKey} class="w-full bg-sidebar border border-border rounded-lg pl-4 pr-12 py-3 text-sm font-mono text-main" />
          <button onclick={handleCopy} class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-muted hover:text-primary transition-colors">
            <Copy size={16} />
          </button>
        </div>
        <div class="pt-2 flex justify-end">
          <button onclick={onClose} class="px-5 py-2 bg-element border border-border text-main text-sm font-bold rounded-lg hover:bg-element-hover transition-colors">Done</button>
        </div>
      </div>
    {/if}
  </div>
</div>
