<script lang="ts">
  let { balanceKobo } = $props<{ balanceKobo: number }>();
  let formattedBalance = $derived(`₦${(balanceKobo / 100).toLocaleString()}`);
</script>

<div class="bg-[var(--bg-sidebar)] border border-[var(--border-color)] rounded-xl p-6">
  <p class="text-[13px] font-medium text-muted uppercase tracking-[0.05em] mb-1 m-0">
    Wallet Balance
  </p>
  <p class="text-3xl font-bold text-main m-0">{formattedBalance}</p>
</div>
