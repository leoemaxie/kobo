<script lang="ts">
  import { ArrowDownLeft, ArrowUpRight, Clock, CheckCircle, AlertCircle } from '@lucide/svelte';

  interface Transaction {
    id: string;
    amount_kobo: number;
    direction: 'inbound' | 'outbound';
    status: 'matched' | 'partial' | 'overpayment' | 'pending' | 'failed';
    monnify_reference: string;
    occurred_at: string;
  }

  let { transactions }: { transactions: Transaction[] } = $props();

  const formatCurrency = (kobo: number) => {
    return new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN' }).format(
      kobo / 100,
    );
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };
</script>

<div class="bg-carbon border border-iron rounded-xl overflow-hidden shadow-sm flex flex-col h-full">
  <div class="p-6 border-b border-iron/50">
    <h3 class="text-sm font-semibold text-smoke uppercase tracking-widest">Transaction History</h3>
  </div>

  <div class="overflow-x-auto">
    <table class="min-w-full">
      <thead>
        <tr class="border-b border-iron/50 bg-void-black/50">
          <th class="px-6 py-3 text-left text-[10px] font-bold uppercase tracking-widest text-smoke"
            >Date</th
          >
          <th class="px-6 py-3 text-left text-[10px] font-bold uppercase tracking-widest text-smoke"
            >Reference</th
          >
          <th class="px-6 py-3 text-left text-[10px] font-bold uppercase tracking-widest text-smoke"
            >Amount</th
          >
          <th
            class="px-6 py-3 text-right text-[10px] font-bold uppercase tracking-widest text-smoke"
            >Status</th
          >
        </tr>
      </thead>
      <tbody class="divide-y divide-iron/30">
        {#each transactions as tx}
          <tr class="hover:bg-graphite/20 transition-colors">
            <td class="px-6 py-3 text-xs text-smoke whitespace-nowrap"
              >{formatDate(tx.occurred_at)}</td
            >
            <td class="px-6 py-3 text-xs font-mono text-smoke whitespace-nowrap"
              >{tx.monnify_reference || tx.id.slice(0, 8)}</td
            >
            <td class="px-6 py-3 whitespace-nowrap">
              <div
                class="flex items-center gap-1.5 text-sm font-medium {tx.direction === 'inbound'
                  ? 'text-electric-lime'
                  : 'text-pure-white'}"
              >
                {#if tx.direction === 'inbound'}
                  <ArrowDownLeft size={14} />
                  + {formatCurrency(tx.amount_kobo)}
                {:else}
                  <ArrowUpRight size={14} class="text-smoke" />
                  - {formatCurrency(tx.amount_kobo)}
                {/if}
              </div>
            </td>
            <td class="px-6 py-3 whitespace-nowrap text-right">
              {#if tx.status === 'matched'}
                <span
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium bg-dark-olive/30 text-electric-lime border border-electric-lime/20"
                >
                  <CheckCircle size={10} /> Matched
                </span>
              {:else if tx.status === 'pending'}
                <span
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium bg-yellow-500/10 text-yellow-500 border border-yellow-500/20"
                >
                  <Clock size={10} /> Pending
                </span>
              {:else}
                <span
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium bg-iron/50 text-paper border border-iron"
                >
                  <AlertCircle size={10} />
                  {tx.status}
                </span>
              {/if}
            </td>
          </tr>
        {/each}
        {#if transactions.length === 0}
          <tr>
            <td colspan="4" class="px-6 py-8 text-center text-sm text-smoke italic">
              No transactions found for this account.
            </td>
          </tr>
        {/if}
      </tbody>
    </table>
  </div>
</div>
