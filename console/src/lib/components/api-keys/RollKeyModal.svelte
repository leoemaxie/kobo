<script lang="ts">
  import { Copy, AlertTriangle, KeySquare } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  let { keyId, onClose }: { keyId: string; onClose: () => void } = $props();

  let secretKey = $state('');
  let isLoading = $state(false);
  let confirmText = $state('');

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }
</script>

<Modal title="Roll API Key" {onClose}>
  {#if !secretKey}
    <form
      method="POST"
      action="?/rollKey"
      use:enhance={() => {
        isLoading = true;
        return async ({ result, update }) => {
          isLoading = false;
          if (result.type === 'success' && (result as any).data?.plainSecret) {
            secretKey = (result as any).data.plainSecret as string;
            toast.success('Key rolled successfully');
          } else {
            toast.error((result as any).data?.error || 'Failed to roll key');
          }
          await update({ reset: false });
        };
      }}
      class="p-6 flex flex-col gap-6"
    >
      <input type="hidden" name="keyId" value={keyId} />

      <div class="border border-error/30 bg-error-bg rounded-lg p-4 flex items-start gap-3">
        <AlertTriangle class="text-error shrink-0 mt-0.5" size={16} />
        <div>
          <h4 class="text-[13px] font-semibold text-error mb-1">Destructive Action</h4>
          <p class="text-xs text-main opacity-90 leading-relaxed">
            Rolling this key will instantly invalidate the current secret. Any systems
            authenticating with <code
              class="font-mono text-[10px] bg-background border border-border px-1 py-0.5 rounded text-main"
              >{keyId}</code
            > will be permanently denied access.
          </p>
        </div>
      </div>

      <div class="space-y-2">
        <label for="confirm" class="block text-xs font-medium text-muted">
          To verify, type <span class="font-bold text-main">roll my key</span> below:
        </label>
        <input
          type="text"
          id="confirm"
          bind:value={confirmText}
          class="w-full bg-sidebar border border-border rounded-lg px-3 py-2.5 text-sm text-main placeholder:text-muted/50 focus:outline-none focus:border-error focus:ring-1 focus:ring-error/20 transition-all"
          placeholder="roll my key"
          autocomplete="off"
        />
      </div>

      <div class="pt-2 flex justify-end gap-3">
        <Button type="button" variant="ghost" size="md" onclick={onClose}>Cancel</Button>
        <Button
          type="submit"
          variant="danger"
          size="md"
          disabled={isLoading || confirmText !== 'roll my key'}
        >
          {isLoading ? 'Rolling...' : 'Roll Key'}
        </Button>
      </div>
    </form>
  {:else}
    <div class="p-6 flex flex-col gap-6">
      <div
        class="bg-primary-transparent border border-primary-border rounded-lg p-4 flex items-start gap-3"
      >
        <KeySquare class="text-primary shrink-0 mt-0.5" size={16} />
        <div>
          <h4 class="text-[13px] font-semibold text-primary mb-1">New Secret Key Generated</h4>
          <p class="text-xs text-main opacity-90 leading-relaxed">
            For security reasons, this is the only time it will be shown. Please copy it and update
            your environment variables immediately.
          </p>
        </div>
      </div>

      <div class="relative group">
        <input
          type="text"
          readonly
          value={secretKey}
          class="w-full bg-sidebar border border-border rounded-lg pl-4 pr-12 py-3 text-sm font-mono text-main outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all selection:bg-primary/20"
        />
        <button
          onclick={handleCopy}
          class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-muted hover:text-primary transition-colors bg-sidebar group-hover:bg-sidebar"
        >
          <Copy size={16} />
        </button>
      </div>

      <div class="pt-2 flex justify-end">
        <Button onclick={onClose} variant="ghost" size="md">I've copied the new key</Button>
      </div>
    </div>
  {/if}
</Modal>
