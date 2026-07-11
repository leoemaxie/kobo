<script lang="ts">
  let { data } = $props();
  let student = $derived(data.student);

  function formatCurrency(kobo: number | string | undefined) {
    if (kobo === undefined || kobo === null) return '₦0.00';
    const num = typeof kobo === 'string' ? parseInt(kobo.replace(/[^0-9]/g, '')) : kobo;
    return new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN' }).format(num / 100);
  }
  
  function formatDate(isoString: string) {
    return new Date(isoString).toLocaleDateString('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit'
    });
  }
</script>

<svelte:head>
  <title>{student.name} | Triumph Academy</title>
</svelte:head>

<div class="space-y-8 w-full">
  <div class="flex items-center gap-4 border-b border-iron pb-6">
    <a href="/dashboard" class="w-10 h-10 rounded-full bg-carbon border border-iron flex items-center justify-center text-smoke hover:text-electric-lime hover:border-electric-lime transition-colors">
      ←
    </a>
    <div>
      <h1 class="text-3xl font-bold text-pure-white tracking-tight">{student.name}</h1>
      <p class="text-smoke mt-1">{student.class} • Account Statement</p>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
    <!-- Left Column: Payment Details -->
    <div class="lg:col-span-1 space-y-6">
      <div class="bg-carbon border border-iron rounded-xl p-6">
        <p class="text-xs text-smoke font-semibold uppercase tracking-widest mb-2">Current Balance</p>
        <h2 class="text-4xl font-bold text-pure-white font-mono tracking-tight">{formatCurrency(student.statement.balance)}</h2>
        
        <div class="mt-6 pt-6 border-t border-iron space-y-4">
          <h3 class="text-sm font-semibold text-pure-white">Make a Payment</h3>
          <p class="text-xs text-smoke leading-relaxed">
            Transfer funds to the virtual account below to instantly credit this student's balance. Powered by Kobo API.
          </p>
          
          <div class="bg-void-black rounded-lg p-4 border border-iron space-y-3">
            <div>
              <p class="text-[10px] text-smoke uppercase font-bold tracking-widest">Bank Name</p>
              <p class="text-sm font-medium text-pure-white mt-0.5">{student.virtualAccount.bankName}</p>
            </div>
            <div>
              <p class="text-[10px] text-smoke uppercase font-bold tracking-widest">Account Number</p>
              <div class="flex items-center justify-between mt-0.5">
                <p class="text-lg font-mono font-bold text-electric-lime tracking-widest">{student.virtualAccount.accountNumber}</p>
                <button class="text-xs bg-carbon px-2 py-1 rounded border border-iron text-smoke hover:text-pure-white transition-colors">Copy</button>
              </div>
            </div>
            <div>
              <p class="text-[10px] text-smoke uppercase font-bold tracking-widest">Account Name</p>
              <p class="text-sm font-medium text-pure-white mt-0.5">{student.virtualAccount.accountName}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Column: Transaction History -->
    <div class="lg:col-span-2">
      <h3 class="text-lg font-semibold text-pure-white mb-4">Transaction History</h3>
      <div class="bg-carbon border border-iron rounded-xl overflow-hidden">
        {#if student.transactions.length > 0}
          <div class="divide-y divide-iron">
            {#each student.transactions as tx}
              <div class="p-4 sm:p-6 flex items-center justify-between hover:bg-graphite/30 transition-colors">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 rounded-full bg-void-black border border-iron flex items-center justify-center">
                    <span class="text-electric-lime font-bold">↓</span>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-pure-white">{tx.description}</p>
                    <p class="text-xs text-smoke mt-1">{formatDate(tx.date)}</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="text-sm font-bold font-mono text-electric-lime">{formatCurrency(tx.amount_kobo ?? tx.amount)}</p>
                  <p class="text-[10px] text-smoke uppercase font-bold tracking-widest mt-1">{tx.status}</p>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <div class="p-12 text-center text-smoke">
            No transactions found for this account.
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
