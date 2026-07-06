<script lang="ts">
  import { Mail, ArrowRight } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import AuthLogo from '$lib/components/ui/AuthLogo.svelte';
  import IconInput from '$lib/components/ui/IconInput.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  let email = $state('');
  let loading = $state(false);
  let submitted = $state(false);
</script>

<div class="w-full">
  <AuthLogo heading="Reset password">
    Remembered it?
    <a href="/auth/login" class="font-semibold text-primary hover:opacity-80 transition-colors ml-1">Sign in</a>
  </AuthLogo>

  <!-- Card -->
  <div class="bg-element border border-border rounded-[10px] px-16 py-8 shadow-sm">
    {#if !submitted}
      <form class="space-y-6" method="POST" use:enhance={() => {
        loading = true;
        return async ({ result, update }) => {
          loading = false;
          if (result.type === 'failure') {
            toast.error(result.data?.error as string || 'Operation failed. Please try again.');
          } else if (result.type === 'error') {
            toast.error('An unexpected server error occurred.');
          } else {
            toast.success('Reset link sent to your email.');
            submitted = true;
          }
          await update();
        };
      }}>
        
        <p class="text-sm text-muted text-center">
          Enter the email associated with your account and we'll send you a link to reset your password.
        </p>

        <IconInput id="email" label="Email Address" type="email" name="email" placeholder="you@company.com" bind:value={email} required>
          {#snippet icon()}<Mail size={15} class="text-subtle" />{/snippet}
        </IconInput>

        <div class="pt-2">
          <Button type="submit" variant="primary" size="lg" class="w-full" disabled={loading}>
            Send Reset Link <ArrowRight size={16} />
          </Button>
        </div>
      </form>
    {:else}
      <div class="text-center space-y-6 py-4">
        <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-background border border-border">
          <Mail size={24} class="text-primary" />
        </div>
        <h2 class="text-lg font-bold text-main">Check your email</h2>
        <p class="text-sm text-muted px-4">
          We sent a password reset link to <span class="text-main font-medium">{email}</span>
        </p>
        <div class="pt-4">
          <button
            type="button"
            onclick={() => submitted = false}
            class="text-sm font-medium text-primary hover:opacity-80 transition-colors"
          >
            Didn't receive it? Click to resend.
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>
