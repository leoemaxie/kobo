<script lang="ts">
  import CurrentPeriodCard from './CurrentPeriodCard.svelte';
  import PlanDetailsCard from './PlanDetailsCard.svelte';
  import InvoiceHistory from './InvoiceHistory.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import CodeBadge from '$lib/components/ui/CodeBadge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';
  import { toast } from '$lib/state/toast.svelte';
  import { page } from '$app/stores';
  import { enhance } from '$app/forms';
  import { onMount } from 'svelte';

  const state = useConsoleState();

  onMount(() => {
    if ($page.url.searchParams.get('payment_success')) {
      toast.success('Payment method successfully added.');
    }
  });
</script>

<svelte:head>
  <title>Billing — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="Billing &amp; Usage">
    {#snippet meta()}
      <span class="font-inconsolata text-[13px] text-subtle">plan:</span>
      <CodeBadge variant="neutral">{state.billingOverview?.plan || 'pay_as_you_go'}</CodeBadge>
      <span class="text-[11px] text-muted">·</span>
      <span class="font-inconsolata text-[13px] text-subtle">
        next_invoice: <span class="text-main">{state.billingOverview?.nextInvoiceDate || '2026-11-01'}</span>
      </span>
    {/snippet}
    {#snippet actions()}
      <form method="POST" action="?/setupPaymentMethod" use:enhance={() => {
        return async ({ result, update }) => {
          if (result.type === 'failure') toast.error((result.data?.error as string) || 'Failed to initialize checkout');
          else if (result.type === 'error') toast.error('Server error occurred');
          await update();
        };
      }}>
        <Button variant="neutral" size="md" type="submit">
          Manage Payment Method
        </Button>
      </form>
    {/snippet}
  </PageHeader>

  <!-- Metrics + plan side-by-side -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-px bg-border-subtle border border-border-subtle rounded-lg overflow-hidden">
    <CurrentPeriodCard />
    <PlanDetailsCard />
  </div>

  <InvoiceHistory />
</div>
