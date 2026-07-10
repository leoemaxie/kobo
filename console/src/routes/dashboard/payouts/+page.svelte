<script lang="ts">
  import type { PageData } from './$types';
  import { Plus, Building2 } from '@lucide/svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import BankAccountModal from '$lib/components/payouts/BankAccountModal.svelte';
  import RequestPayoutModal from '$lib/components/payouts/RequestPayoutModal.svelte';
  import WalletBalanceCard from '$lib/components/payouts/WalletBalanceCard.svelte';
  import BankAccountCard from '$lib/components/payouts/BankAccountCard.svelte';
  import PayoutsHistory from '$lib/components/payouts/PayoutsHistory.svelte';

  let { data } = $props<{ data: PageData }>();

  let isBankModalOpen = $state(false);
  let isPayoutModalOpen = $state(false);

  let hasBankAccount = $derived(!!data.bankAccount);
</script>

<svelte:head>
  <title>Payouts — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="Payouts">
    {#snippet meta()}
      <span class="text-[13px] text-muted"
        >Manage your settlement bank account and request payouts.</span
      >
    {/snippet}
    {#snippet actions()}
      <Button variant="neutral" size="md" onclick={() => (isBankModalOpen = true)}>
        <Building2 size={16} />
        {hasBankAccount ? 'Update Bank' : 'Add Bank'}
      </Button>
      <Button
        variant="primary"
        size="md"
        onclick={() => (isPayoutModalOpen = true)}
        disabled={!hasBankAccount}
        title={!hasBankAccount ? 'Please add a bank account first' : ''}
      >
        <Plus size={16} />
        Request Payout
      </Button>
    {/snippet}
  </PageHeader>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <WalletBalanceCard balanceKobo={data.walletBalanceKobo} />
    <BankAccountCard bankAccount={data.bankAccount} onConfigure={() => (isBankModalOpen = true)} />
  </div>

  <PayoutsHistory payouts={data.payouts} />
</div>

<BankAccountModal
  isOpen={isBankModalOpen}
  onClose={() => (isBankModalOpen = false)}
  banks={data.banks || []}
/>

<RequestPayoutModal
  isOpen={isPayoutModalOpen}
  onClose={() => (isPayoutModalOpen = false)}
  walletBalanceKobo={data.walletBalanceKobo}
/>
