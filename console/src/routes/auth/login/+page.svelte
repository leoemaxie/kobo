<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Mail, Lock, Eye, EyeOff, LogIn } from '@lucide/svelte';

  let email = $state('');
  let password = $state('');
  let showPassword = $state(false);
  let loading = $state(false);
</script>

<div class="w-full">
  <!-- Logo + heading -->
  <div class="mb-8 text-center">
    <!-- Use an invert filter in light mode for the logo if it's white, or keep it if we have a theme-aware logo -->
    <img src="/logo.png" alt="Kobo" class="h-10 w-auto mx-auto mb-6" style="filter: var(--logo-filter);" />
    <h1 class="text-2xl font-inter font-bold text-main tracking-tight">Sign in to Console</h1>
    <p class="text-muted text-sm mt-2 leading-relaxed">Enter your credentials to access the console</p>
  </div>

  <!-- Card -->
  <div class="bg-element border border-border rounded-[10px] px-16 py-6 shadow-sm">
    <form class="space-y-6" method="POST">
      
      <!-- Email field -->
      <div class="space-y-1.5">
        <label for="email" class="block text-xs font-semibold text-muted uppercase tracking-widest">Email</label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
            <Mail size={15} class="text-subtle" />
          </div>
          <input
            id="email"
            type="email"
            bind:value={email}
            required
            placeholder="you@company.com"
            class="block w-full rounded-[8px] border border-border bg-background pl-10 pr-4 py-3 text-sm text-main placeholder-subtle focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors"
          />
        </div>
      </div>

      <!-- Password field -->
      <div class="space-y-1.5">
        <div class="flex items-center justify-between">
          <label for="password" class="block text-xs font-semibold text-muted uppercase tracking-widest">Password</label>
          <a href="/auth/forgot-password" class="text-xs font-medium text-primary hover:opacity-80 transition-colors">Forgot password?</a>
        </div>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
            <Lock size={15} class="text-subtle" />
          </div>
          <input
            id="password"
            type={showPassword ? 'text' : 'password'}
            bind:value={password}
            required
            placeholder="••••••••••••"
            class="block w-full rounded-[8px] border border-border bg-background pl-10 pr-10 py-3 text-sm text-main placeholder-subtle focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors"
          />
          <button
            type="button"
            onclick={() => (showPassword = !showPassword)}
            class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-subtle hover:text-muted transition-colors"
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

      <!-- Remember me -->
      <div class="flex items-center gap-3">
        <input 
          id="remember" 
          type="checkbox" 
          class="h-5 w-5 flex-shrink-0 appearance-none rounded-[4px] border border-border bg-element checked:bg-primary checked:border-primary focus:outline-none focus:ring-1 focus:ring-primary focus:ring-offset-1 focus:ring-offset-background cursor-pointer transition-colors relative before:content-[''] before:absolute before:inset-0 before:bg-[url('data:image/svg+xml;utf8,%3Csvg%20viewBox%3D%220%200%2016%2016%22%20fill%3D%22%23080808%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%3Cpath%20d%3D%22M12.207%204.793a1%201%200%20010%201.414l-5%205a1%201%200%2001-1.414%200l-2-2a1%201%200%20011.414-1.414L6.5%209.086l4.293-4.293a1%201%200%20011.414%200z%22%2F%3E%3C%2Fsvg%3E')] before:bg-no-repeat before:bg-center before:bg-[length:14px_14px] before:opacity-0 checked:before:opacity-100" 
        />
        <label for="remember" class="text-sm text-muted cursor-pointer select-none">Keep me signed in for 7 days</label>
      </div>

      <div class="pt-1">
        <button
          type="submit"
          disabled={loading}
          class="w-full flex items-center justify-center gap-2 rounded-[8px] bg-primary text-primary-text px-6 py-2.5 text-sm font-bold tracking-tight shadow-md hover:opacity-90 hover:-translate-y-0.5 active:translate-y-0 transition-all focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 focus:ring-offset-element disabled:opacity-60 disabled:cursor-not-allowed"
        >
          <LogIn size={16} />
          Sign In
        </button>
      </div>
    </form>
  </div>

  <!-- Footer -->
  <p class="text-center text-sm text-muted mt-6">
    No account yet?
    <a href="/auth/signup" class="font-semibold text-primary hover:opacity-80 transition-colors ml-1">Sign up</a>
  </p>
</div>
