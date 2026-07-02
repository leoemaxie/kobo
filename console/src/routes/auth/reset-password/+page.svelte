<script lang="ts">
  import { Lock, Eye, EyeOff, Save } from '@lucide/svelte';

  let password = $state('');
  let confirmPassword = $state('');
  let showPassword = $state(false);
  let showConfirmPassword = $state(false);
  let loading = $state(false);
</script>

<div class="w-full">
  <!-- Logo + heading -->
  <div class="mb-8 text-center">
    <img src="/logo.png" alt="Kobo" class="h-10 w-auto mx-auto mb-6" />
    <h1 class="text-3xl font-inter font-bold text-pure-white tracking-tight">Set new password</h1>
    <p class="text-smoke text-sm mt-3 leading-relaxed">
      Must be at least 8 characters long
    </p>
  </div>

  <!-- Card -->
  <div class="bg-carbon border border-iron rounded-[10px] px-16 py-8 shadow-sm">
    <form class="space-y-6" method="POST">
      
      <!-- New Password field -->
      <div class="space-y-1.5">
        <label for="password" class="block text-xs font-semibold text-smoke uppercase tracking-widest">New Password</label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
            <Lock size={15} class="text-fog" />
          </div>
          <input
            id="password"
            type={showPassword ? 'text' : 'password'}
            bind:value={password}
            required
            placeholder="••••••••••••"
            class="block w-full rounded-[8px] border border-iron bg-void-black pl-10 pr-10 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
          />
          <button
            type="button"
            onclick={() => (showPassword = !showPassword)}
            class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-fog hover:text-smoke transition-colors"
            aria-label={showPassword ? 'Hide password' : 'Show password'}
          >
            {#if showPassword}
              <EyeOff size={15} />
            {:else}
              <Eye size={15} />
            {/if}
          </button>
        </div>
      </div>

      <!-- Confirm Password field -->
      <div class="space-y-1.5">
        <label for="confirmPassword" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Confirm Password</label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
            <Lock size={15} class="text-fog" />
          </div>
          <input
            id="confirmPassword"
            type={showConfirmPassword ? 'text' : 'password'}
            bind:value={confirmPassword}
            required
            placeholder="••••••••••••"
            class="block w-full rounded-[8px] border border-iron bg-void-black pl-10 pr-10 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
          />
          <button
            type="button"
            onclick={() => (showConfirmPassword = !showConfirmPassword)}
            class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-fog hover:text-smoke transition-colors"
            aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
          >
            {#if showConfirmPassword}
              <EyeOff size={15} />
            {:else}
              <Eye size={15} />
            {/if}
          </button>
        </div>
      </div>

      <div class="pt-2">
        <button
          type="submit"
          disabled={loading || (password.length > 0 && password !== confirmPassword)}
          class="w-full flex items-center justify-center gap-2 rounded-[8px] bg-electric-lime text-void-black px-6 py-2.5 text-sm font-bold tracking-tight shadow-md hover:bg-lime-glow hover:-translate-y-0.5 active:translate-y-0 transition-all focus:outline-none focus:ring-2 focus:ring-electric-lime focus:ring-offset-2 focus:ring-offset-carbon disabled:opacity-60 disabled:cursor-not-allowed"
        >
          <Save size={16} />
          Update Password
        </button>
      </div>
      
      {#if password.length > 0 && password !== confirmPassword && confirmPassword.length > 0}
        <p class="text-xs text-red-400 text-center mt-2">Passwords do not match.</p>
      {/if}
    </form>
  </div>
</div>
