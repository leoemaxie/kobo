<script lang="ts">
  import { enhance } from '$app/forms';
  import { X } from '@lucide/svelte';
  import Select, { type SelectOption } from '$lib/components/ui/Select.svelte';

  interface Bank {
    name: string;
    code: string;
  }

  let { isOpen, onClose, banks } = $props<{
    isOpen: boolean;
    onClose: () => void;
    banks: Bank[];
  }>();

  let loading = $state(false);
  let selectedBank = $state('');
  let accountNumber = $state('');
  let submitError = $state('');

  $effect(() => {
    // Automatically strip non-digits and enforce 10 chars
    const cleaned = accountNumber.replace(/\D/g, '').slice(0, 10);
    if (accountNumber !== cleaned) {
      accountNumber = cleaned;
    }
  });
  let bankName = $derived(banks.find((b: Bank) => b.code === selectedBank)?.name || '');
  let bankOptions = $derived(
    banks.map((b: Bank): SelectOption => ({ value: b.code, label: b.name })),
  );
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-0">
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 bg-[var(--bg-sidebar)] opacity-80" onclick={onClose}></div>
    <div
      class="relative bg-[var(--bg-element)] border border-[var(--border-color)] rounded-xl shadow-lg w-full max-w-md p-6"
    >
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-lg font-semibold text-main m-0">Set Bank Account</h2>
        <button
          class="bg-transparent border-none cursor-pointer p-1 text-muted hover:text-main"
          onclick={onClose}
        >
          <X size={20} />
        </button>
      </div>

      {#if submitError}
        <div
          class="mb-4 p-3 rounded-lg bg-[var(--error-bg)] border border-[var(--error-color)] text-[var(--error-color)] text-[13px] font-medium"
        >
          {submitError}
        </div>
      {/if}

      <form
        method="POST"
        action="?/saveBankAccount"
        class="flex flex-col gap-4"
        use:enhance={({ cancel }) => {
          submitError = '';

          if (!selectedBank) {
            submitError = 'Please select a bank from the list.';
            cancel();
            return;
          }

          if (accountNumber.length !== 10) {
            submitError = 'Account number must be exactly 10 digits.';
            cancel();
            return;
          }

          loading = true;
          return async ({ result, update }) => {
            loading = false;
            if (result.type === 'success') {
              onClose();
            } else if (result.type === 'failure') {
              submitError = result.data?.error?.toString() || 'Invalid bank account details.';
            }
            update();
          };
        }}
      >
        <div class="flex flex-col gap-1.5">
          <label for="bankCode" class="text-[13px] font-medium text-main">Bank Name</label>
          <Select
            id="bankCode"
            name="bankCode"
            bind:value={selectedBank}
            options={bankOptions}
            placeholder="Select a bank"
          /><input type="hidden" name="bankName" value={bankName} />
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="accountNumber" class="text-[13px] font-medium text-main">Account Number</label
          >
          <input
            type="text"
            id="accountNumber"
            name="accountNumber"
            required
            minlength="10"
            maxlength="10"
            placeholder="0123456789"
            bind:value={accountNumber}
            class="w-full h-9 px-3 rounded-lg border border-[var(--border-color)] bg-transparent text-[13px] text-main outline-none focus:border-[var(--accent)] placeholder:text-muted"
          />
        </div>

        <div class="mt-2 flex justify-end gap-3">
          <button
            type="button"
            class="h-9 px-4 rounded-lg bg-[var(--bg-element)] text-main text-[13px] font-medium border border-[var(--border-color)] cursor-pointer"
            onclick={onClose}
            disabled={loading}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="h-9 px-4 rounded-lg bg-[var(--accent)] text-[var(--accent-text)] text-[13px] font-semibold border-none cursor-pointer disabled:opacity-50"
            disabled={loading}
          >
            {loading ? 'Verifying...' : 'Save Account'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
