<script lang="ts">
  import { Landmark } from '@lucide/svelte';
  let { payouts } = $props<{ payouts: any[] }>();

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  function getStatusColor(status: string) {
    switch (status.toLowerCase()) {
      case 'successful':
        return 'bg-[#10b981]/10 text-[#10b981] border-[#10b981]/20';
      case 'failed':
        return 'bg-[var(--error-bg)] text-[var(--error-color)] border-[var(--error-color)]';
      case 'processing':
        return 'bg-[#f59e0b]/10 text-[#f59e0b] border-[#f59e0b]/20';
      default:
        return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  }
</script>

<div class="bg-[var(--bg-sidebar)] border border-[var(--border-color)] rounded-xl overflow-hidden">
  <div class="p-5 border-b border-[var(--border-color)]">
    <h3 class="text-[15px] font-semibold text-main m-0">Recent Payouts</h3>
  </div>

  {#if payouts.length > 0}
    <div class="overflow-x-auto">
      <table class="w-full border-collapse">
        <thead>
          <tr class="border-b border-[var(--border-color)] bg-[var(--bg-element)]">
            <th
              class="px-5 py-3 text-left text-[11px] font-semibold text-muted uppercase tracking-[0.05em]"
              >Amount</th
            >
            <th
              class="px-5 py-3 text-left text-[11px] font-semibold text-muted uppercase tracking-[0.05em]"
              >Status</th
            >
            <th
              class="px-5 py-3 text-left text-[11px] font-semibold text-muted uppercase tracking-[0.05em]"
              >Date</th
            >
            <th
              class="px-5 py-3 text-right text-[11px] font-semibold text-muted uppercase tracking-[0.05em]"
              >Ref</th
            >
          </tr>
        </thead>
        <tbody>
          {#each payouts as payout}
            <tr
              class="border-b border-[var(--border-color)] last:border-0 hover:bg-[var(--bg-element)] transition-colors"
            >
              <td class="px-5 py-4 whitespace-nowrap text-[14px] font-semibold text-main">
                ₦{(payout.requested_amount_kobo / 100).toLocaleString()}
              </td>
              <td class="px-5 py-4 whitespace-nowrap">
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-md text-[11px] font-medium border uppercase tracking-wider {getStatusColor(
                    payout.status,
                  )}"
                >
                  {payout.status}
                </span>
              </td>
              <td class="px-5 py-4 whitespace-nowrap text-[13px] text-muted">
                {formatDate(payout.created_at)}
              </td>
              <td class="px-5 py-4 whitespace-nowrap text-[13px] text-muted text-right font-mono">
                {payout.merchant_tx_ref}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else}
    <div class="p-12 text-center">
      <div
        class="w-12 h-12 rounded-full bg-[var(--bg-element)] border border-[var(--border-color)] flex items-center justify-center mx-auto mb-3 text-muted"
      >
        <Landmark size={24} />
      </div>
      <p class="text-[14px] font-medium text-main m-0">No payouts yet</p>
      <p class="text-[13px] text-muted mt-1 max-w-sm mx-auto">
        Your payout history will appear here once you request a withdrawal.
      </p>
    </div>
  {/if}
</div>
