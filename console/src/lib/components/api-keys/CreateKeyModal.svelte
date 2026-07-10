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

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }
</script>

<Modal title="Create Standard Key" {onClose}>
  {#if !secretKey}
    <form method="POST" action="?/createKey" use:enhance={() => {
      isLoading = true;
      return async ({ result, update }) => {
        isLoading = false;
        if (result.type === 'success' && (result as any).data?.plainSecret) {
          secretKey = (result as any).data.plainSecret as string;
          toast.success('Key created successfully');
        } else {
          toast.error((result as any).data?.error || 'Failed to create key');
        }
        await update({ reset: false });
      };
    }} class="p-6 flex flex-col gap-5">
      <Input id="label" label="Key Label" type="text" name="label" placeholder="e.g. Production Main" required />
      
      <div class="space-y-1.5">
        <label for="environment" class="block text-xs font-semibold text-muted uppercase tracking-widest">Environment</label>
        <input type="hidden" name="environment" value={currentEnv} />
        <div class="w-full bg-sidebar border border-border rounded-lg px-4 py-2.5 text-sm text-main capitalize opacity-70 cursor-not-allowed">
          {currentEnv}
        </div>
      </div>

      <div class="pt-2 flex justify-end gap-3">
        <Button type="button" variant="ghost" size="md" onclick={onClose}>Cancel</Button>
        <Button type="submit" variant="primary" size="md" disabled={isLoading}>
          {isLoading ? 'Creating...' : 'Create Key'}
        </Button>
      </div>
    </form>
  {:else}
    <div class="p-6 flex flex-col gap-5">
      <div class="bg-primary-transparent border border-primary-border rounded-lg p-4">
        <p class="text-sm font-medium text-primary mb-2">Please copy this secret key now.</p>
        <p class="text-xs text-muted leading-relaxed">For security reasons, this is the only time it will be shown. If you lose it, you will need to roll the key.</p>
      </div>

      <div class="relative">
        <input type="text" readonly value={secretKey} class="w-full bg-sidebar border border-border rounded-lg pl-4 pr-12 py-3 text-sm font-mono text-main" />
        <button onclick={handleCopy} class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-muted hover:text-primary transition-colors">
          <Copy size={16} />
        </button>
      </div>

      <div class="pt-2 flex justify-end">
        <Button onclick={onClose} variant="ghost" size="md">I've copied the key</Button>
      </div>
    </div>
  {/if}
</Modal>
