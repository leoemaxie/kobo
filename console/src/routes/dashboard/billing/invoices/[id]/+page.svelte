<script lang="ts">
  import type { PageData } from './$types';
  import { onMount } from 'svelte';

  export let data: PageData;
  const { invoice, integrator, details } = data;

  onMount(() => {
    // Automatically open print dialog when loading this specific route
    // setTimeout(() => window.print(), 500);
  });
</script>

<svelte:head>
  <title>Invoice {invoice.id} - Kobo Console</title>
  <style>
    @media print {
      body {
        background: white !important;
        color: black !important;
      }
    }
  </style>
</svelte:head>

<div
  class="max-w-4xl my-10 mx-auto p-4 sm:p-10 border border-border-subtle bg-element rounded-lg overflow-x-auto print:border-[#ddd] print:m-0 print:p-0 print:border-none print:shadow-none"
>
  <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start mb-10 gap-4 sm:gap-0">
    <div>
      <h1 class="text-2xl font-bold m-0 text-main print:text-black">KOBO INC.</h1>
      <p class="text-sm mt-1 text-subtle print:text-black">Receipt of Payment</p>
    </div>
    <div class="text-left sm:text-right">
      <h2 class="text-xl font-bold m-0 text-main print:text-black">INVOICE</h2>
      <p class="font-mono text-sm mt-1 text-subtle print:text-black">
        {invoice.id}
      </p>
    </div>
  </div>

  <div class="flex flex-col sm:flex-row sm:justify-between mb-10 gap-4 sm:gap-0">
    <div>
      <p class="text-xs uppercase tracking-widest text-subtle mb-2 print:text-black">Billed To</p>
      <p class="text-base font-medium text-main m-0 print:text-black">
        {integrator}
      </p>
    </div>
    <div class="text-left sm:text-right">
      <p class="text-xs uppercase tracking-widest text-subtle mb-2 print:text-black">
        Date of Issue
      </p>
      <p class="text-base font-medium text-main m-0 print:text-black">
        {invoice.date}
      </p>
    </div>
  </div>

  <table class="w-full border-collapse mb-10">
    <thead>
      <tr class="border-b border-border-subtle text-left print:border-[#ddd]">
        <th class="py-3 text-xs uppercase tracking-widest text-subtle print:text-black"
          >Description</th
        >
        <th class="py-3 text-xs uppercase tracking-widest text-subtle text-right print:text-black"
          >Amount</th
        >
      </tr>
    </thead>
    <tbody>
      <tr class="border-b border-border-subtle print:border-[#ddd]">
        <td class="py-4 text-main print:text-black">Platform Usage - {invoice.period}</td>
        <td class="py-4 font-mono text-right text-main print:text-black">{invoice.amount}</td>
      </tr>
      {#if details}
        <tr class="border-b border-border-subtle print:border-[#ddd]">
          <td class="py-4 text-[13px] text-subtle print:text-black">
            Includes {details.accounts} virtual accounts, {details.transactions}
            transactions, {details.webhooks} webhooks.
          </td>
          <td class="py-4 print:text-black"></td>
        </tr>
      {/if}
    </tbody>
    <tfoot>
      <tr>
        <td class="pt-6 pb-2 font-bold text-main print:text-black">Total Due</td>
        <td class="pt-6 pb-2 font-mono text-lg font-bold text-right text-main print:text-black"
          >{invoice.amount}</td
        >
      </tr>
    </tfoot>
  </table>

  <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 sm:gap-0">
    <div>
      <span
        class="text-xs font-bold uppercase py-1 px-2 rounded {invoice.status === 'paid'
          ? 'bg-primary-transparent text-primary'
          : 'bg-error-bg text-error'}"
      >
        {invoice.status}
      </span>
    </div>

    <button
      class="print:hidden border border-border rounded-md bg-sidebar py-2 px-4 text-[13px] font-semibold text-muted cursor-pointer hover:bg-element-hover transition-colors"
      on:click={() => window.print()}
    >
      Download PDF
    </button>
  </div>
</div>
