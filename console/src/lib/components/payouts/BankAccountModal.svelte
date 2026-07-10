<script lang="ts">
	import { enhance } from '$app/forms';
	import { X } from '@lucide/svelte';

	let { isOpen, onClose, banks } = $props<{ 
		isOpen: boolean, 
		onClose: () => void,
		banks: Array<{ name: string, code: string }>
	}>();

	let loading = $state(false);
	let selectedBank = $state('');
	
	// Automatically find bankName based on bankCode
	let bankName = $derived(banks.find(b => b.code === selectedBank)?.name || '');
</script>

{#if isOpen}
<div class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-0">
	<div class="fixed inset-0 bg-[var(--bg-sidebar)] opacity-80" onclick={onClose}></div>
	<div class="relative bg-[var(--bg-main)] border border-[var(--border-color)] rounded-xl shadow-lg w-full max-w-md p-6 overflow-hidden">
		<div class="flex justify-between items-center mb-6">
			<h2 class="text-lg font-semibold text-main m-0">Set Bank Account</h2>
			<button class="bg-transparent border-none cursor-pointer p-1 text-muted hover:text-main" onclick={onClose}>
				<X size={20} />
			</button>
		</div>

		<form method="POST" action="?/saveBankAccount" class="flex flex-col gap-4" use:enhance={() => {
			loading = true;
			return async ({ result, update }) => {
				loading = false;
				if (result.type === 'success') {
					onClose();
				}
				update();
			};
		}}>
			<div class="flex flex-col gap-1.5">
				<label for="bankCode" class="text-[13px] font-medium text-main">Bank Name</label>
				<select 
					id="bankCode" 
					name="bankCode" 
					bind:value={selectedBank}
					required 
					class="w-full h-9 px-3 rounded-lg border border-[var(--border-color)] bg-transparent text-[13px] text-main outline-none focus:border-[var(--accent)]"
				>
					<option value="" disabled>Select a bank</option>
					{#each banks as bank}
						<option value={bank.code}>{bank.name}</option>
					{/each}
				</select>
				<input type="hidden" name="bankName" value={bankName} />
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="accountNumber" class="text-[13px] font-medium text-main">Account Number</label>
				<input 
					type="text" 
					id="accountNumber" 
					name="accountNumber" 
					required 
					pattern="[0-9]{10}"
					title="Must be 10 digits"
					maxlength="10"
					placeholder="0123456789"
					class="w-full h-9 px-3 rounded-lg border border-[var(--border-color)] bg-transparent text-[13px] text-main outline-none focus:border-[var(--accent)] placeholder:text-muted"
				/>
			</div>

			<div class="mt-2 flex justify-end gap-3">
				<button type="button" class="h-9 px-4 rounded-lg bg-[var(--bg-element)] text-main text-[13px] font-medium border border-[var(--border-color)] cursor-pointer" onclick={onClose} disabled={loading}>
					Cancel
				</button>
				<button type="submit" class="h-9 px-4 rounded-lg bg-[var(--accent)] text-[var(--accent-text)] text-[13px] font-semibold border-none cursor-pointer disabled:opacity-50" disabled={loading}>
					{loading ? 'Verifying...' : 'Save Account'}
				</button>
			</div>
		</form>
	</div>
</div>
{/if}
