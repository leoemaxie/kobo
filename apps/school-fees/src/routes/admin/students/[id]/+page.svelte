<script lang="ts">
  import PageHeader from '$lib/components/common/PageHeader.svelte';
  import FinancialSummary from '$lib/components/admin/students/details/FinancialSummary.svelte';
  import LinkedParents from '$lib/components/admin/students/details/LinkedParents.svelte';
  import TransactionHistory from '$lib/components/admin/students/details/TransactionHistory.svelte';
  import StudentActions from '$lib/components/admin/students/details/StudentActions.svelte';
  import { ArrowLeft } from '@lucide/svelte';

  let { data, form } = $props();
  let student = $derived(data.student);
  let availableParents = $derived(data.availableParents);
</script>

<svelte:head>
  <title>{student.name} | Admin Console</title>
</svelte:head>

<div class="space-y-8 w-full pb-10">
  <!-- Navigation and Header -->
  <div>
    <a
      href="/admin/students"
      class="inline-flex items-center gap-2 text-sm text-smoke hover:text-pure-white transition-colors mb-6"
    >
      <ArrowLeft size={16} />
      Back to Students
    </a>

    <PageHeader
      title={student.name}
      subtitle="Student ID: {student.id} &bull; Class: {student.class}"
    />
  </div>

  {#if form?.success}
    <div
      class="bg-dark-olive/20 border border-electric-lime/50 text-electric-lime text-sm p-4 rounded-lg shadow-sm"
    >
      {form.message || 'Operation successful'}
    </div>
  {/if}
  {#if form?.error}
    <div class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg shadow-sm">
      {form.error}
    </div>
  {/if}

  <!-- Top Section: Summaries and Actions -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
    <FinancialSummary account={student.virtualAccount} statement={student.statement} />
    <LinkedParents linkedParents={student.linkedParents} {availableParents} />
    <StudentActions studentId={student.id} currentClass={student.class} />
  </div>

  <!-- Bottom Section: Transactions -->
  <div class="w-full">
    <TransactionHistory transactions={student.transactions} />
  </div>
</div>
