<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { CreditCard, Download, FileText, ArrowUpRight } from 'lucide-svelte';
</script>

<div class="space-y-8">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-inter font-bold text-pure-white mb-2 tracking-tight">Billing & Usage</h1>
      <p class="text-bone">View your current usage and manage billing details.</p>
    </div>
    <Button variant="ghost" class="gap-2">
      <CreditCard size={16} /> Payment Methods
    </Button>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
    <Card>
      <h2 class="text-sm font-semibold tracking-widest text-smoke uppercase mb-6">Current Billing Period</h2>
      <div class="flex items-end gap-3 mb-2">
        <span class="text-5xl font-inter font-bold text-pure-white">₦14,250</span>
        <span class="text-smoke text-sm mb-1 pb-1">/ Oct 1 - Oct 31, 2026</span>
      </div>
      
      <div class="mt-8 space-y-4">
        <div>
          <div class="flex justify-between text-sm mb-1">
            <span class="text-paper">Virtual Accounts Provisioned</span>
            <span class="text-bone font-mono">152 × ₦50 = ₦7,600</span>
          </div>
          <div class="w-full bg-graphite h-1.5 rounded-full overflow-hidden">
            <div class="bg-electric-lime h-full w-[45%]"></div>
          </div>
        </div>
        
        <div>
          <div class="flex justify-between text-sm mb-1">
            <span class="text-paper">Transaction Fees</span>
            <span class="text-bone font-mono">3,325 × ₦2 = ₦6,650</span>
          </div>
          <div class="w-full bg-graphite h-1.5 rounded-full overflow-hidden">
            <div class="bg-electric-lime h-full w-[35%] opacity-75"></div>
          </div>
        </div>
      </div>
    </Card>

    <Card class="flex flex-col justify-between">
      <div>
        <h2 class="text-sm font-semibold tracking-widest text-smoke uppercase mb-4">Plan Details</h2>
        <h3 class="text-2xl font-basier font-semibold text-pure-white">Pay as you go</h3>
        <p class="text-smoke text-sm mt-2">
          You are currently on the flexible Pay As You Go plan. Upgrade to a committed volume tier for lower per-transaction fees.
        </p>
      </div>
      <div class="mt-8">
        <Button variant="primary" class="w-full sm:w-auto">View Upgrade Options</Button>
      </div>
    </Card>
  </div>

  <div class="space-y-4">
    <h2 class="text-xl font-basier font-semibold text-pure-white mt-4">Invoice History</h2>
    <Card class="p-0 overflow-hidden">
      <table class="min-w-full divide-y divide-iron">
        <thead class="bg-carbon">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Date</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Amount</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Status</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Invoice</th>
          </tr>
        </thead>
        <tbody class="bg-void-black divide-y divide-iron">
          {#each [
            { date: 'Sep 1, 2026', amount: '₦32,400', status: 'Paid' },
            { date: 'Aug 1, 2026', amount: '₦28,150', status: 'Paid' },
            { date: 'Jul 1, 2026', amount: '₦24,900', status: 'Paid' }
          ] as invoice}
            <tr class="hover:bg-graphite/50 transition-colors">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-paper">{invoice.date}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-bone font-mono">{invoice.amount}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex items-center rounded-radius-full bg-dark-olive px-2.5 py-0.5 text-xs font-medium text-electric-lime">
                  {invoice.status}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button class="text-smoke hover:text-electric-lime transition-colors flex items-center gap-1">
                  <Download size={14} /> Download PDF
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </Card>
  </div>
</div>
