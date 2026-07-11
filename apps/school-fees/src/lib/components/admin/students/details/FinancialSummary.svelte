<script lang="ts">
  import { Copy, Check, CreditCard, Landmark } from '@lucide/svelte';

  interface VirtualAccount {
    accountName: string | null;
    accountNumber: string | null;
    bankName: string;
  }

  interface Statement {
    balanceKobo: number;
    currency: string;
  }

  let { account, statement }: { account: VirtualAccount; statement: Statement } = $props();
  
  let copied = $state(false);

  function copyAccount() {
    if (account.accountNumber) {
      navigator.clipboard.writeText(account.accountNumber);
      copied = true;
      setTimeout(() => copied = false, 2000);
    }
  }

  const formatCurrency = (kobo: number) => {
    return new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN' }).format(kobo / 100);
  };
</script>

<div class="bg-carbon border border-iron rounded-xl p-6 h-full shadow-sm flex flex-col gap-8">
  <div>
    <h3 class="text-sm font-semibold text-smoke uppercase tracking-widest mb-4">Financial Summary</h3>
    
    <div class="">
      <div class="text-xs text-smoke mb-1">Current Balance</div>
      <div class="text-4xl font-bold text-pure-white mt-2">{formatCurrency(statement.balanceKobo)}</div>
    </div>
  </div>

  <div class="bg-void-black border border-iron/50 rounded-lg p-4">
    <div class="flex items-center gap-2 mb-3">
      <Landmark size={16} class="text-electric-lime" />
      <span class="text-xs font-semibold text-smoke uppercase tracking-widest">Virtual Account</span>
    </div>
    
    {#if account.accountNumber}
      <div class="space-y-3">
        <div>
          <div class="text-[10px] text-smoke uppercase tracking-wider">Account Number</div>
          <div class="flex items-center gap-2 mt-0.5">
            <span class="text-lg font-medium text-electric-lime">{account.accountNumber}</span>
            <button onclick={copyAccount} class="text-smoke hover:text-electric-lime focus:outline-none transition-colors">
              {#if copied}
                <Check size={16} class="text-electric-lime" />
              {:else}
                <Copy size={16} />
              {/if}
            </button>
          </div>
        </div>
        
        <div class="grid grid-cols-2 gap-4">
          <div>
            <div class="text-[10px] text-smoke uppercase tracking-wider">Bank</div>
            <div class="text-sm font-medium text-pure-white mt-0.5">{account.bankName}</div>
          </div>
          <div>
            <div class="text-[10px] text-smoke uppercase tracking-wider">Name</div>
            <div class="text-sm font-medium text-pure-white mt-0.5 truncate">{account.accountName}</div>
          </div>
        </div>
      </div>
    {:else}
      <div class="text-sm text-smoke italic py-4">Account provisioning pending or failed.</div>
    {/if}
  </div>
</div>
