<script lang="ts">
  import { enhance } from '$app/forms';
  import { X, Info } from '@lucide/svelte';

  let {
    isOpen,
    onClose,
    walletBalanceKobo,
    minPayoutKobo = 100000,
    transferFeeBufferKobo = 5000,
  } = $props<{
    isOpen: boolean;
    onClose: () => void;
    walletBalanceKobo: number;
    minPayoutKobo?: number;
    transferFeeBufferKobo?: number;
  }>();

  let loading = $state(false);
  let amountString = $state('');

  let amountKobo = $derived(Math.floor(parseFloat(amountString || '0') * 100));

  let maxAvailableKobo = $derived(Math.max(0, walletBalanceKobo - transferFeeBufferKobo));

  let isAmountValid = $derived(amountKobo >= minPayoutKobo && amountKobo <= maxAvailableKobo);
  let isAmountTooHigh = $derived(amountKobo > maxAvailableKobo);
  let isAmountTooLow = $derived(amountKobo > 0 && amountKobo < minPayoutKobo);

  function setMax() {
    amountString = (maxAvailableKobo / 100).toString();
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-0">
    <div class="fixed inset-0 bg-[var(--bg-sidebar)] opacity-80" onclick={onClose}></div>
    <div
      class="relative bg-[var(--bg-main)] border border-[var(--border-color)] rounded-xl shadow-lg w-full max-w-md p-6 overflow-hidden"
    >
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-lg font-semibold text-main m-0">Request Payout</h2>
        <button
          class="bg-transparent border-none cursor-pointer p-1 text-muted hover:text-main"
          onclick={onClose}
        >
          <X size={20} />
        </button>
      </div>

      <form
        method="POST"
        action="?/requestPayout"
        class="flex flex-col gap-4"
        use:enhance={() => {
          loading = true;
          return async ({ result, update }) => {
            loading = false;
            if (result.type === 'success') {
              onClose();
            }
            update();
          };
        }}
      >
        <div
          class="p-3 bg-[var(--bg-element)] border border-[var(--border-color)] rounded-lg flex gap-3"
        >
          <Info size={18} class="text-muted mt-0.5" />
          <div>
            <p class="text-[13px] text-main font-medium m-0 mb-1">
              Available for Payout: ₦{(maxAvailableKobo / 100).toLocaleString()}
            </p>
            <p class="text-[12px] text-muted m-0">
              ₦{(transferFeeBufferKobo / 100).toLocaleString()} is reserved for Nomba transfer fees.
            </p>
          </div>
        </div>

        <div class="flex flex-col gap-1.5 mt-2">
          <label
            for="amount"
            class="flex items-center justify-between text-[13px] font-medium text-main"
          >
            <span>Amount (₦)</span>
            <button
              type="button"
              class="text-[12px] text-[var(--accent)] bg-transparent border-none p-0 cursor-pointer hover:underline"
              onclick={setMax}>Max Available</button
            >
          </label>
          <input
            type="number"
            id="amount"
            name="amount"
            bind:value={amountString}
            required
            min={minPayoutKobo / 100}
            max={maxAvailableKobo / 100}
            step="0.01"
            placeholder="e.g. 5000"
            class="w-full h-9 px-3 rounded-lg border border-[var(--border-color)] bg-transparent text-[13px] text-main outline-none focus:border-[var(--accent)] placeholder:text-muted"
            class:border-[var(--error-color)]={isAmountTooHigh || isAmountTooLow}
          />
          {#if isAmountTooHigh}
            <p class="text-[11px] text-[var(--error-color)] m-0">
              Amount exceeds maximum available (₦{(maxAvailableKobo / 100).toLocaleString()}).
            </p>
          {:else if isAmountTooLow}
            <p class="text-[11px] text-[var(--error-color)] m-0">
              Minimum payout is ₦{(minPayoutKobo / 100).toLocaleString()}.
            </p>
          {/if}
        </div>

        <div class="mt-4 flex justify-end gap-3">
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
            disabled={loading || !isAmountValid}
          >
            {loading ? 'Processing...' : 'Request Payout'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
