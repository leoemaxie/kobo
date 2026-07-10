<script lang="ts">
  import { Copy } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let { onClose }: { onClose: () => void } = $props();

  const consoleState = useConsoleState();
  let currentEnv = $derived(consoleState.currentEnvironment);

  let secretKey = $state('');
  let isLoading = $state(false);
  let selectedEvents: Record<string, boolean> = $state({
    'account.created': false,
    'transaction.completed': false,
    'transaction.failed': false,
    'account.closed': false,
  });

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }

  let eventsList = $derived(
    Object.keys(selectedEvents)
      .filter((k) => selectedEvents[k])
      .join(','),
  );
</script>

<Modal title="Add Webhook Endpoint" {onClose}>
  {#if !secretKey}
    <form
      method="POST"
      action="?/addEndpoint"
      use:enhance={() => {
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
      }}
      class="p-6 flex flex-col gap-5"
    >
      <input type="hidden" name="environment" value={currentEnv} />

      <Input
        id="url"
        label="Endpoint URL"
        type="url"
        name="url"
        placeholder="https://api.yourdomain.com/webhooks"
        required
        class="font-mono"
      />

      <div class="space-y-1.5">
        <label for="events" class="block text-xs font-semibold text-muted uppercase tracking-widest"
          >Events to Listen For</label
        >
        <input id="events" type="hidden" name="events" value={eventsList} />
        <div class="flex flex-col gap-2 bg-sidebar border border-border p-3 rounded-lg">
          {#each Object.keys(selectedEvents) as event}
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                bind:checked={selectedEvents[event]}
                class="rounded bg-element border-border text-primary focus:ring-primary focus:ring-offset-sidebar"
              />
              <span class="text-sm font-mono text-main">{event}</span>
            </label>
          {/each}
        </div>
      </div>

      <div class="pt-2 flex justify-end gap-3">
        <Button type="button" variant="ghost" size="md" onclick={onClose}>Cancel</Button>
        <Button type="submit" variant="primary" size="md" disabled={isLoading || !eventsList}>
          {isLoading ? 'Adding...' : 'Add Endpoint'}
        </Button>
      </div>
    </form>
  {:else}
    <div class="p-6 flex flex-col gap-5">
      <div class="bg-primary-transparent border border-primary-border rounded-lg p-4">
        <p class="text-sm font-medium text-primary mb-2">Please copy this webhook secret now.</p>
        <p class="text-xs text-muted leading-relaxed">
          Use this secret to verify webhook payloads. It will not be shown again.
        </p>
      </div>
      <div class="relative">
        <input
          type="text"
          readonly
          value={secretKey}
          class="w-full bg-sidebar border border-border rounded-lg pl-4 pr-12 py-3 text-sm font-mono text-main"
        />
        <button
          onclick={handleCopy}
          class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-muted hover:text-primary transition-colors"
        >
          <Copy size={16} />
        </button>
      </div>
      <div class="pt-2 flex justify-end">
        <Button onclick={onClose} variant="ghost" size="md">Done</Button>
      </div>
    </div>
  {/if}
</Modal>
