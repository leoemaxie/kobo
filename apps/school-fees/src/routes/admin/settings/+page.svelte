<script lang="ts">
  import { enhance } from '$app/forms';
  let { form } = $props();
  let isSubmitting = $state(false);
</script>

<svelte:head>
  <title>Settings | Triumph Academy</title>
</svelte:head>

<div class="space-y-10 w-full">
  <div class="flex items-center justify-between border-b border-iron pb-6">
    <div>
      <h1 class="text-3xl font-bold text-pure-white tracking-tight">Kobo Integration Settings</h1>
      <p class="text-smoke mt-1">Configure short-lived API credentials for interacting with Kobo services.</p>
    </div>
  </div>

  <div class="bg-carbon border border-iron rounded-xl p-8 max-w-2xl">
    <div class="mb-6 space-y-2">
      <div class="bg-electric-lime/10 border border-electric-lime/30 text-electric-lime text-sm p-4 rounded-lg">
        <strong>Note:</strong> These credentials are kept in memory and are not saved to the database. They will reset when the server restarts.
      </div>
    </div>
    
    {#if form?.error}
      <div class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg mb-6 shadow-sm">
        {form.error}
      </div>
    {/if}
    {#if form?.success}
      <div class="bg-dark-olive/20 border border-electric-lime/50 text-electric-lime text-sm p-4 rounded-lg mb-6 shadow-sm">
        Credentials updated in memory successfully!
      </div>
    {/if}

    <form method="POST" class="space-y-6" use:enhance={() => {
      isSubmitting = true;
      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}>
      <div class="space-y-1.5">
        <label for="apiKey" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Kobo API Key</label>
        <input
          type="text"
          id="apiKey"
          name="apiKey"
          required
          placeholder="pk_test_..."
          class="block w-full rounded-lg border border-iron bg-void-black px-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
        />
      </div>

      <div class="space-y-1.5">
        <label for="apiSecret" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Kobo API Secret</label>
        <input
          type="password"
          id="apiSecret"
          name="apiSecret"
          required
          placeholder="sk_test_..."
          class="block w-full rounded-lg border border-iron bg-void-black px-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
        />
      </div>

      <div class="pt-4">
        <button
          type="submit"
          disabled={isSubmitting}
          class="rounded-lg bg-electric-lime text-void-black px-6 py-3 text-sm font-bold shadow-md hover:bg-lime-glow transition-all disabled:opacity-50"
        >
          {isSubmitting ? 'Saving...' : 'Set Credentials'}
        </button>
      </div>
    </form>
  </div>
</div>
