<script lang="ts">
  import { X, Copy, AlertTriangle } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';

  let { keyId, onClose }: { keyId: string, onClose: () => void } = $props();

  let secretKey = $state('');
  let isLoading = $state(false);

  function handleCopy() {
    navigator.clipboard.writeText(secretKey);
    toast.success('Secret copied to clipboard');
  }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
  <div class="w-full max-w-md bg-element border border-border rounded-xl overflow-hidden shadow-xl">
    <div class="flex items-center justify-between px-6 py-4 border-b border-border bg-sidebar">
      <h2 class="text-sm font-semibold text-main uppercase tracking-widest">Roll API Key</h2>
      <button onclick={onClose} class="text-muted hover:text-main transition-colors">
        <X size={16} />
      </button>
    </div>

    {#if !secretKey}
      <form method="POST" action="?/rollKey" use:enhance={() => {
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
      }} class="p-6 flex flex-col gap-5">
        <input type="hidden" name="keyId" value={keyId} />
        
        <div class="bg-error-bg border border-error/30 rounded-lg p-4 flex gap-3">
          <AlertTriangle class="text-error shrink-0 mt-0.5" size={18} />
          <div>
            <p class="text-sm font-medium text-error mb-1">This will invalidate the current key.</p>
            <p class="text-xs text-muted leading-relaxed">Any applications currently using <code class="font-mono text-[10px] bg-background px-1 py-0.5 rounded text-main">{keyId}</code> will stop working immediately. A new key will be generated with the same label and permissions.</p>
          </div>
        </div>

        <div class="pt-2 flex justify-end gap-3">
          <button type="button" onclick={onClose} class="px-4 py-2 text-sm font-medium text-muted hover:text-main transition-colors">Cancel</button>
          <button type="submit" disabled={isLoading} class="px-5 py-2 bg-error text-white text-sm font-bold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50">
            {isLoading ? 'Rolling...' : 'Roll Key'}
          </button>
        </div>
      </form>
    {:else}
      <div class="p-6 flex flex-col gap-5">
        <div class="bg-primary-transparent border border-primary-border rounded-lg p-4">
          <p class="text-sm font-medium text-primary mb-2">Please copy your NEW secret key.</p>
          <p class="text-xs text-muted leading-relaxed">Update your applications with this new secret immediately.</p>
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
