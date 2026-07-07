<script lang="ts">
  import { Mail, Lock, Eye, EyeOff, LogIn } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import AuthLogo from '$lib/components/ui/AuthLogo.svelte';
  import IconInput from '$lib/components/ui/IconInput.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  let email = $state('');
  let password = $state('');
  let showPassword = $state(false);
  let loading = $state(false);
</script>

<svelte:head>
  <title>Login | Kobo Console</title>
</svelte:head>


<div class="w-full">
  <AuthLogo heading="Sign in to Console">
    <span>Enter your credentials to access the console</span>
  </AuthLogo>

  <!-- Card -->
  <div class="bg-element border border-border rounded-[10px] px-16 py-6 shadow-sm">
    <form class="space-y-6" method="POST" use:enhance={() => {
      loading = true;
      return async ({ result, update }) => {
        loading = false;
        if (result.type === 'failure') {
          toast.error(result.data?.error as string || 'Login failed. Please check your credentials.');
        } else if (result.type === 'error') {
          toast.error('An unexpected server error occurred.');
        }
        await update();
      };
    }}>
      
      <!-- Email field -->
      <IconInput id="email" label="Email" type="email" name="email" placeholder="you@company.com" bind:value={email} required>
        {#snippet icon()}<Mail size={15} class="text-subtle" />{/snippet}
      </IconInput>

      <!-- Password field -->
      <div class="space-y-1.5">
        <div class="flex items-center justify-between">
          <label for="password" class="block text-xs font-semibold text-muted uppercase tracking-widest">Password</label>
          <a href="/auth/forgot-password" class="text-xs font-medium text-primary hover:opacity-80 transition-colors">Forgot password?</a>
        </div>
        <IconInput id="password" label="" type={showPassword ? 'text' : 'password'} name="password" placeholder="••••••••••••" bind:value={password} required>
          {#snippet icon()}<Lock size={15} class="text-subtle" />{/snippet}
          {#snippet trailing()}
            <button
              type="button"
              onclick={() => (showPassword = !showPassword)}
              class="pr-3.5 flex items-center text-subtle hover:text-muted transition-colors"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
            >
              {#if showPassword}<EyeOff size={15} />{:else}<Eye size={15} />{/if}
            </button>
          {/snippet}
        </IconInput>
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
        <Button type="submit" variant="primary" size="lg" class="w-full" disabled={loading}>
          <LogIn size={16} /> Sign In
        </Button>
      </div>
    </form>
  </div>

  <!-- Footer -->
  <p class="text-center text-sm text-muted mt-6">
    No account yet?
    <a href="/auth/signup" class="font-semibold text-primary hover:opacity-80 transition-colors ml-1">Sign up</a>
  </p>
</div>
