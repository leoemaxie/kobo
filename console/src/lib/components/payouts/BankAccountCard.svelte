<script lang="ts">
  import { Landmark } from '@lucide/svelte';
  let { bankAccount, onConfigure } = $props<{
    bankAccount: any;
    onConfigure: () => void;
  }>();
  let hasBankAccount = $derived(!!bankAccount);
</script>

<div
  class="bg-[var(--bg-sidebar)] border border-[var(--border-color)] rounded-xl p-6 flex flex-col justify-center"
>
  {#if hasBankAccount}
    <p class="text-[13px] font-medium text-muted uppercase tracking-[0.05em] mb-3 m-0">
      Settlement Account
    </p>
    <div class="flex items-center gap-4">
      <div
        class="w-12 h-12 bg-[var(--bg-element)] border border-[var(--border-color)] rounded-full flex items-center justify-center text-muted"
      >
        <Landmark size={24} />
      </div>
      <div>
        <p class="text-[15px] font-semibold text-main m-0">{bankAccount.account_name}</p>
        <p class="text-[13px] text-muted m-0">
          {bankAccount.account_number} • {bankAccount.bank_name}
        </p>
      </div>
    </div>
  {:else}
    <div class="text-center">
      <p class="text-[13px] text-muted mb-3">No bank account configured.</p>
      <button
        class="text-[13px] font-medium text-[var(--accent)] bg-transparent border-none p-0 cursor-pointer hover:underline"
        onclick={onConfigure}
      >
        Configure one now &rarr;
      </button>
    </div>
  {/if}
</div>
