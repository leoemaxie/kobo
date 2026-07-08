<script lang="ts">
  import { Lock, Eye, EyeOff, Save } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import AuthLogo from '$lib/components/ui/AuthLogo.svelte';
  import IconInput from '$lib/components/ui/IconInput.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  let { data } = $props();

  let password = $state('');
  let confirmPassword = $state('');
  let showPassword = $state(false);
  let showConfirmPassword = $state(false);
  let loading = $state(false);

  let mismatch = $derived(password.length > 0 && confirmPassword.length > 0 && password !== confirmPassword);
</script>

<svelte:head>
  <title>Reset Password | Kobo Console</title>
</svelte:head>


<div class="w-full">
  <AuthLogo heading="Set new password">
    <span>Must be at least 8 characters long</span>
  </AuthLogo>

  <!-- Card -->
  <div class="bg-element border border-border rounded-[10px] px-16 py-8 shadow-sm">
    <form class="space-y-6" method="POST" use:enhance={() => {
      loading = true;
      return async ({ result, update }) => {
        loading = false;
        if (result.type === 'failure') {
          toast.error(result.data?.error as string || 'Password reset failed. Please try again.');
        } else if (result.type === 'error') {
          toast.error('An unexpected server error occurred.');
        } else if (result.type === 'redirect' || result.type === 'success') {
          toast.success('Password updated successfully!');
        }
        await update();
      };
    }}>
      <input type="hidden" name="token" value={data.token} />

      <!-- New Password -->
      <IconInput id="password" label="New Password" type={showPassword ? 'text' : 'password'} name="password" placeholder="••••••••••••" bind:value={password} required>
        {#snippet icon()}<Lock size={15} class="text-subtle" />{/snippet}
        {#snippet trailing()}
          <button type="button" onclick={() => (showPassword = !showPassword)} class="pr-3.5 flex items-center text-subtle hover:text-muted transition-colors" aria-label={showPassword ? 'Hide password' : 'Show password'}>
            {#if showPassword}<EyeOff size={15} />{:else}<Eye size={15} />{/if}
          </button>
        {/snippet}
      </IconInput>

      <!-- Confirm Password -->
      <IconInput id="confirmPassword" label="Confirm Password" type={showConfirmPassword ? 'text' : 'password'} name="confirmPassword" placeholder="••••••••••••" bind:value={confirmPassword} required>
        {#snippet icon()}<Lock size={15} class="text-subtle" />{/snippet}
        {#snippet trailing()}
          <button type="button" onclick={() => (showConfirmPassword = !showConfirmPassword)} class="pr-3.5 flex items-center text-subtle hover:text-muted transition-colors" aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}>
            {#if showConfirmPassword}<EyeOff size={15} />{:else}<Eye size={15} />{/if}
          </button>
        {/snippet}
      </IconInput>

      {#if mismatch}
        <p class="text-xs text-red-500 text-center">Passwords do not match.</p>
      {/if}

      <div class="pt-2">
        <Button type="submit" variant="primary" size="lg" class="w-full" disabled={loading || mismatch}>
          <Save size={16} /> Update Password
        </Button>
      </div>
    </form>
  </div>
</div>
