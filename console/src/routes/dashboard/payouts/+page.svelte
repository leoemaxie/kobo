<script lang="ts">
	import type { PageData } from './$types';
	import { Plus, Landmark, Building2, ExternalLink } from '@lucide/svelte';
	import BankAccountModal from '$lib/components/payouts/BankAccountModal.svelte';
	import RequestPayoutModal from '$lib/components/payouts/RequestPayoutModal.svelte';
	
	let { data } = $props<{ data: PageData }>();

	let isBankModalOpen = $state(false);
	let isPayoutModalOpen = $state(false);

	let formattedBalance = $derived(`₦${(data.walletBalanceKobo / 100).toLocaleString()}`);
	let hasBankAccount = $derived(!!data.bankAccount);

	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleDateString(undefined, {
			year: 'numeric', month: 'short', day: 'numeric',
			hour: '2-digit', minute: '2-digit'
		});
	}

	function getStatusColor(status: string) {
		switch (status.toLowerCase()) {
			case 'successful': return 'bg-[#10b981]/10 text-[#10b981] border-[#10b981]/20';
			case 'failed': return 'bg-[var(--error-bg)] text-[var(--error-color)] border-[var(--error-color)]';
			case 'processing': return 'bg-[#f59e0b]/10 text-[#f59e0b] border-[#f59e0b]/20';
			default: return 'bg-gray-100 text-gray-700 border-gray-200';
		}
	}
</script>

<div class="flex-1 overflow-y-auto min-h-0 bg-[var(--bg-main)]">
	<div class="max-w-6xl mx-auto p-4 sm:p-8">
		<header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
			<div>
				<h1 class="text-2xl font-bold text-main m-0 mb-1">Payouts</h1>
				<p class="text-[13px] text-muted m-0">Manage your settlement bank account and request payouts.</p>
			</div>
			
			<div class="flex gap-3">
				<button 
					class="flex items-center gap-2 h-9 px-4 rounded-lg bg-[var(--bg-element)] text-main text-[13px] font-medium border border-[var(--border-color)] cursor-pointer hover:bg-[var(--bg-active)]"
					onclick={() => isBankModalOpen = true}
				>
					<Building2 size={16} />
					{hasBankAccount ? 'Update Bank' : 'Add Bank'}
				</button>
				<button 
					class="flex items-center gap-2 h-9 px-4 rounded-lg bg-[var(--accent)] text-[var(--accent-text)] text-[13px] font-semibold border-none cursor-pointer hover:opacity-90 disabled:opacity-50"
					onclick={() => isPayoutModalOpen = true}
					disabled={!hasBankAccount}
					title={!hasBankAccount ? 'Please add a bank account first' : ''}
				>
					<Plus size={16} />
					Request Payout
				</button>
			</div>
		</header>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
			<!-- Wallet Balance Card -->
			<div class="bg-[var(--bg-sidebar)] border border-[var(--border-color)] rounded-xl p-6">
				<p class="text-[13px] font-medium text-muted uppercase tracking-[0.05em] mb-1 m-0">Wallet Balance</p>
				<p class="text-3xl font-bold text-main m-0">{formattedBalance}</p>
			</div>

			<!-- Bank Account Card -->
			<div class="bg-[var(--bg-sidebar)] border border-[var(--border-color)] rounded-xl p-6 flex flex-col justify-center">
				{#if hasBankAccount}
					<p class="text-[13px] font-medium text-muted uppercase tracking-[0.05em] mb-3 m-0">Settlement Account</p>
					<div class="flex items-center gap-4">
						<div class="w-12 h-12 bg-[var(--bg-element)] border border-[var(--border-color)] rounded-full flex items-center justify-center text-muted">
							<Landmark size={24} />
						</div>
						<div>
							<p class="text-[15px] font-semibold text-main m-0">{data.bankAccount.accountName}</p>
							<p class="text-[13px] text-muted m-0">{data.bankAccount.accountNumber} • {data.bankAccount.bankName}</p>
						</div>
					</div>
				{:else}
					<div class="text-center">
						<p class="text-[13px] text-muted mb-3">No bank account configured.</p>
						<button class="text-[13px] font-medium text-[var(--accent)] bg-transparent border-none p-0 cursor-pointer hover:underline" onclick={() => isBankModalOpen = true}>
							Configure one now &rarr;
						</button>
					</div>
				{/if}
			</div>
		</div>

		<!-- Payouts History -->
		<div class="bg-[var(--bg-sidebar)] border border-[var(--border-color)] rounded-xl overflow-hidden">
			<div class="p-5 border-b border-[var(--border-color)]">
				<h3 class="text-[15px] font-semibold text-main m-0">Recent Payouts</h3>
			</div>

			{#if data.payouts.length > 0}
				<div class="overflow-x-auto">
					<table class="w-full border-collapse">
						<thead>
							<tr class="border-b border-[var(--border-color)] bg-[var(--bg-element)]">
								<th class="px-5 py-3 text-left text-[11px] font-semibold text-muted uppercase tracking-[0.05em]">Amount</th>
								<th class="px-5 py-3 text-left text-[11px] font-semibold text-muted uppercase tracking-[0.05em]">Status</th>
								<th class="px-5 py-3 text-left text-[11px] font-semibold text-muted uppercase tracking-[0.05em]">Date</th>
								<th class="px-5 py-3 text-right text-[11px] font-semibold text-muted uppercase tracking-[0.05em]">Ref</th>
							</tr>
						</thead>
						<tbody>
							{#each data.payouts as payout}
								<tr class="border-b border-[var(--border-color)] last:border-0 hover:bg-[var(--bg-element)] transition-colors">
									<td class="px-5 py-4 whitespace-nowrap text-[14px] font-semibold text-main">
										₦{(payout.requested_amount_kobo / 100).toLocaleString()}
									</td>
									<td class="px-5 py-4 whitespace-nowrap">
										<span class="inline-flex items-center px-2 py-0.5 rounded-md text-[11px] font-medium border uppercase tracking-wider {getStatusColor(payout.status)}">
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
					<div class="w-12 h-12 rounded-full bg-[var(--bg-element)] border border-[var(--border-color)] flex items-center justify-center mx-auto mb-3 text-muted">
						<Landmark size={24} />
					</div>
					<p class="text-[14px] font-medium text-main m-0">No payouts yet</p>
					<p class="text-[13px] text-muted mt-1 max-w-sm mx-auto">Your payout history will appear here once you request a withdrawal.</p>
				</div>
			{/if}
		</div>
	</div>
</div>

<BankAccountModal 
	isOpen={isBankModalOpen} 
	onClose={() => isBankModalOpen = false} 
	banks={data.banks || []} 
/>

<RequestPayoutModal 
	isOpen={isPayoutModalOpen} 
	onClose={() => isPayoutModalOpen = false} 
	walletBalanceKobo={data.walletBalanceKobo} 
/>
