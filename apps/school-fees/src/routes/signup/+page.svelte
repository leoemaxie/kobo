<script lang="ts">
  import { enhance } from '$app/forms';
  let { form } = $props();

  let name = $state('');
  let email = $state('');
  let password = $state('');
  let isSubmitting = $state(false);
</script>

<div class="w-full max-w-md space-y-8">
  <div class="text-center">
    <div class="w-12 h-12 rounded-lg bg-electric-lime text-void-black flex items-center justify-center text-2xl font-bold mx-auto mb-4 shadow-lg shadow-electric-lime/20">K</div>
    <h1 class="text-3xl font-bold text-pure-white tracking-tight">Create Account</h1>
    <p class="text-smoke mt-2">Register for the Triumph Academy parent portal</p>
  </div>

  <div class="bg-carbon border border-iron rounded-xl p-8 shadow-sm">
    {#if form?.error}
      <div class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg mb-6 shadow-sm">
        {form.error}
      </div>
    {/if}

    <form class="space-y-5" method="POST" use:enhance={() => {
      isSubmitting = true;
      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}>
      <div class="space-y-1.5">
        <label for="name" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Full Name</label>
        <input
          id="name"
          name="name"
          type="text"
          bind:value={name}
          required
          placeholder="Jane Doe"
          class="block w-full rounded-lg border border-iron bg-void-black px-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
        />
      </div>

      <div class="space-y-1.5">
        <label for="email" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Email Address</label>
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
        <label for="role" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Role</label>
        <select
          id="role"
          name="role"
          required
          class="block w-full rounded-lg border border-iron bg-void-black px-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors appearance-none"
        >
          <option value="parent">Parent</option>
          <option value="admin">Administrator</option>
        </select>
      </div>

      <div class="space-y-1.5">
        <label for="password" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Password</label>
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
          {isSubmitting ? 'Registering...' : 'Register'}
        </button>
      </div>
    </form>
  </div>

  <p class="text-center text-sm text-smoke">
    Already have an account? 
    <a href="/login" class="font-medium text-electric-lime hover:text-lime-glow transition-colors">Sign in</a>
  </p>
</div>
