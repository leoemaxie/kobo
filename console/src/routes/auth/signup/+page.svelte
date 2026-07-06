<script lang="ts">
  import { Mail, Lock, Eye, EyeOff, UserPlus } from '@lucide/svelte';
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

<div class="w-full">
  <AuthLogo heading="Create your account">
    Already have an account?
    <a href="/auth/login" class="font-semibold text-primary hover:opacity-80 transition-colors ml-1">Sign in</a>
  </AuthLogo>

  <!-- Card -->
  <div class="bg-element border border-border rounded-[10px] px-16 py-6 shadow-sm">
    <form class="space-y-6" method="POST" use:enhance={() => {
      loading = true;
      return async ({ result, update }) => {
        loading = false;
        if (result.type === 'failure') {
          toast.error(result.data?.error as string || 'Signup failed. Please check your details.');
        } else if (result.type === 'error') {
          toast.error('An unexpected server error occurred.');
        }
        await update();
      };
    }}>

      <!-- Email field -->
      <IconInput id="email" label="Work Email" type="email" name="email" placeholder="you@company.com" bind:value={email} required>
        {#snippet icon()}<Mail size={15} class="text-subtle" />{/snippet}
      </IconInput>

      <!-- Password field -->
      <IconInput id="password" label="Password" type={showPassword ? 'text' : 'password'} name="password" placeholder="••••••••••••" bind:value={password} required>
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

      <div class="pt-1">
        <Button type="submit" variant="primary" size="lg" class="w-full" disabled={loading}>
          <UserPlus size={16} /> Create Account
        </Button>
      </div>

      <p class="text-[11px] text-subtle text-center leading-relaxed">
        By creating an account, you agree to our Terms of Service and Privacy Policy.
      </p>
    </form>
  </div>
</div>
