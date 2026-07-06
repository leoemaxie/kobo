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
			.no-print {
				display: none !important;
			}
			.print-text {
				color: black !important;
			}
			.print-border {
				border-color: #ddd !important;
			}
		}
	</style>
</svelte:head>

<div class="print-border" style="max-width: 800px; margin: 40px auto; padding: clamp(16px, 5vw, 40px); border: 1px solid var(--border-subtle); background: var(--bg-surface); border-radius: 8px; overflow-x: auto;">
	<div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 40px;">
		<div>
			<h1 class="print-text" style="font-size: 24px; margin: 0; color: white;">KOBO INC.</h1>
			<p class="print-text" style="color: var(--text-subtle); margin: 4px 0 0; font-size: 14px;">Receipt of Payment</p>
		</div>
		<div style="text-align: right;">
			<h2 class="print-text" style="font-size: 20px; margin: 0; color: white;">INVOICE</h2>
			<p class="print-text" style="font-family: monospace; color: var(--text-subtle); margin: 4px 0 0;">{invoice.id}</p>
		</div>
	</div>

	<div style="display: flex; justify-content: space-between; margin-bottom: 40px;">
		<div>
			<p class="print-text" style="font-size: 12px; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 8px;">Billed To</p>
			<p class="print-text" style="font-size: 16px; font-weight: 500; color: white; margin: 0;">{integrator}</p>
		</div>
		<div style="text-align: right;">
			<p class="print-text" style="font-size: 12px; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 8px;">Date of Issue</p>
			<p class="print-text" style="font-size: 16px; font-weight: 500; color: white; margin: 0;">{invoice.date}</p>
		</div>
	</div>

	<table style="width: 100%; border-collapse: collapse; margin-bottom: 40px;">
		<thead>
			<tr class="print-border" style="border-bottom: 1px solid var(--border-subtle); text-align: left;">
				<th class="print-text" style="padding: 12px 0; font-size: 12px; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle);">Description</th>
				<th class="print-text" style="padding: 12px 0; font-size: 12px; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle); text-align: right;">Amount</th>
			</tr>
		</thead>
		<tbody>
			<tr class="print-border" style="border-bottom: 1px solid var(--border-subtle);">
				<td class="print-text" style="padding: 16px 0; color: white;">Platform Usage - {invoice.period}</td>
				<td class="print-text" style="padding: 16px 0; font-family: monospace; text-align: right; color: white;">{invoice.amount}</td>
			</tr>
			{#if details}
			<tr class="print-border" style="border-bottom: 1px solid var(--border-subtle);">
				<td class="print-text" style="padding: 16px 0; color: var(--text-subtle); font-size: 13px;">
					Includes {details.accounts} virtual accounts, {details.transactions} transactions, {details.webhooks} webhooks.
				</td>
				<td class="print-text" style="padding: 16px 0;"></td>
			</tr>
			{/if}
		</tbody>
		<tfoot>
			<tr>
				<td class="print-text" style="padding: 24px 0 8px; font-weight: 700; color: white;">Total Due</td>
				<td class="print-text" style="padding: 24px 0 8px; font-family: monospace; font-size: 18px; font-weight: 700; text-align: right; color: white;">{invoice.amount}</td>
			</tr>
		</tfoot>
	</table>

	<div style="display: flex; justify-content: space-between; align-items: center;">
		<div>
			<span style="font-size: 12px; font-weight: 700; text-transform: uppercase; padding: 4px 8px; border-radius: 4px; background: {invoice.status === 'paid' ? 'rgba(0, 255, 0, 0.1)' : 'rgba(255, 0, 0, 0.1)'}; color: {invoice.status === 'paid' ? 'var(--accent)' : '#ff4444'};">
				{invoice.status}
			</span>
		</div>
		
		<button class="no-print" on:click={() => window.print()} style="
			border: 1px solid #2a2a2a; border-radius: 6px;
			background: var(--bg-sidebar); padding: 8px 16px;
			font-size: 13px; font-weight: 600; color: var(--text-muted); cursor: pointer;
		">
			Download PDF
		</button>
	</div>
</div>
