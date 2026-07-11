<script lang="ts">
  import { enhance } from '$app/forms';

  let { form } = $props();

  let email = $state('');
  let password = $state('');
  let isSubmitting = $state(false);
</script>

<div class="bg-carbon border border-iron rounded-xl p-8 shadow-sm">
  {#if form?.error}
    <div
      class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg mb-6 shadow-sm"
    >
      {form.error}
    </div>
  {/if}

  <form
    class="space-y-5"
    method="POST"
    use:enhance={() => {
      isSubmitting = true;
      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}
  >
    <div class="space-y-1.5">
      <label for="email" class="block text-xs font-semibold text-smoke uppercase tracking-widest"
        >Email Address</label
      >
      <input
        id="email"
        name="email"
        type="email"
        bind:value={email}
        required
        placeholder="parent@example.com"
        class="block w-full rounded-lg border border-iron bg-void-black px-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
      />
    </div>

    <div class="space-y-1.5">
      <label for="password" class="block text-xs font-semibold text-smoke uppercase tracking-widest"
        >Password</label
      >
      <input
        id="password"
        name="password"
        type="password"
        bind:value={password}
        required
        placeholder="••••••••••••"
        class="block w-full rounded-lg border border-iron bg-void-black px-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
      />
    </div>

    <div class="pt-2">
      <button
        type="submit"
        disabled={isSubmitting}
        class="w-full rounded-lg bg-electric-lime text-void-black px-4 py-3 text-sm font-bold shadow-md hover:bg-lime-glow transition-all disabled:opacity-50"
      >
        {isSubmitting ? 'Signing in...' : 'Sign In'}
      </button>
    </div>
  </form>
</div>
