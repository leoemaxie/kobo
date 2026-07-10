<script lang="ts">
  import { Mail, Check, AlertCircle } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import AuthLogo from '$lib/components/ui/AuthLogo.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  export let data;

  let loading = false;
</script>

<svelte:head>
  <title>Accept Invitation | Kobo Console</title>
</svelte:head>

<div class="w-full">
  <AuthLogo heading="Workspace Invitation" />

  <div class="bg-element border border-border rounded-[10px] px-6 sm:px-12 py-8 shadow-sm">
    {#if data.error}
      <div class="flex flex-col items-center text-center">
        <div
          class="h-12 w-12 rounded-full bg-red-500/10 flex items-center justify-center mb-4 text-red-500"
        >
          <AlertCircle size={24} />
        </div>
        <h3 class="text-lg font-semibold text-main mb-2">Invalid Invitation</h3>
        <p class="text-sm text-subtle mb-6">{data.error}</p>
        <Button href="/auth/login" variant="secondary" size="md">Go to Login</Button>
      </div>
    {:else}
      <div class="flex flex-col items-center text-center">
        <div
          class="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 text-primary"
        >
          <Mail size={24} />
        </div>
        <h3 class="text-lg font-semibold text-main mb-2">You've been invited!</h3>
        <p class="text-[14px] text-muted mb-6 leading-relaxed">
          You have been invited to join the <strong class="text-main font-semibold"
            >{data.invite.workspaceName}</strong
          >
          workspace as a
          <strong class="text-main font-semibold capitalize">{data.invite.role}</strong>.
        </p>

        {#if !data.user}
          <div class="w-full bg-sidebar border border-border rounded-lg p-5 mb-6 text-left">
            <h4 class="text-sm font-medium text-main mb-2">Account Required</h4>
            <p class="text-[13px] text-subtle mb-4">
              You need a Kobo account to accept this invitation.
            </p>
            <div class="flex flex-col gap-3">
              <Button href="/auth/signup" variant="primary" size="md" class="w-full">
                Create Account
              </Button>
              <Button href="/auth/login" variant="secondary" size="md" class="w-full">
                Log In
              </Button>
            </div>
          </div>
        {:else}
          {#if data.isEmailMismatch}
            <div
              class="w-full bg-red-500/10 border border-red-500/20 rounded-lg p-4 mb-6 text-left"
            >
              <p class="text-[13px] text-red-400">
                This invitation was sent to a different email address. Please log out and log in
                with the correct account to accept it.
              </p>
            </div>
          {/if}
          {#if data.user.integratorId}
            <div
              class="w-full bg-orange-500/10 border border-orange-500/20 rounded-lg p-4 mb-6 text-left"
            >
              <p class="text-[13px] text-orange-400">
                You already belong to a workspace. You must leave your current workspace in Settings
                before joining a new one.
              </p>
            </div>
            <Button href="/dashboard" variant="secondary" size="md" class="w-full">
              Go to Dashboard
            </Button>
          {:else}
            <form
              method="POST"
              class="w-full"
              use:enhance={() => {
                loading = true;
                return async ({ result, update }) => {
                  loading = false;
                  if (result.type === 'failure') {
                    toast.error((result.data?.error as string) || 'Failed to accept invitation.');
                  }
                  await update();
                };
              }}
            >
              <input type="hidden" name="token" value={data.invite.id} />
              <Button
                type="submit"
                variant="primary"
                size="lg"
                class="w-full"
                disabled={loading || data.isEmailMismatch}
              >
                <Check size={16} /> Accept Invitation
              </Button>
            </form>
          {/if}
        {/if}
      </div>
    {/if}
  </div>
</div>
