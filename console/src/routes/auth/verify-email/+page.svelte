<script lang="ts">
  import { enhance } from '$app/forms';
  import { Mail, ArrowLeft, RefreshCw, CheckCircle2 } from '@lucide/svelte';
  import { toast } from '$lib/state/toast.svelte';
  import type { PageData } from './$types';
  import Button from '$lib/components/ui/Button.svelte';

  let { data }: { data: PageData } = $props();
  let resending = $state(false);
  let resent = $state(false);
</script>

<svelte:head>
  <title>Verify your email — Kobo Console</title>
</svelte:head>

<div class="w-full text-center">
  <!-- Logo -->
  <div class="mb-10">
    <img src="/logo.png" alt="Kobo" class="h-10 w-auto mx-auto" style="filter: var(--logo-filter);" />
  </div>

  <!-- Main Card -->
  <div class="bg-element border border-border rounded-[10px] px-12 py-10 shadow-sm flex flex-col items-center max-w-[420px] mx-auto">
    
    <!-- Icon -->
    <div class="h-16 w-16 bg-primary/10 rounded-full flex items-center justify-center mb-6">
      <Mail class="text-primary" size={15} />
    </div>

    <!-- Headers -->
    <h1 class="text-2xl font-inter font-bold text-main tracking-tight mb-3">Check your inbox</h1>
    
    <p class="text-sm text-muted leading-relaxed max-w-[280px] mb-8">
      We've sent a verification link to 
      {#if data.email}
        <strong class="text-main font-semibold">{data.email}</strong>.
      {:else}
        your email address.
      {/if}
      <br/><br/>
      Please click the link to activate your account.
    </p>

    <!-- Error state -->
    {#if data?.error}
      <div class="mb-6 w-full text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-[6px] px-4 py-3 text-center">
        {data.error}
      </div>
    {/if}

    <!-- Divider -->
    <div class="w-full border-t border-border pt-6">
      {#if resent}
        <div class="flex flex-col items-center gap-2">
          <CheckCircle2 class="text-primary" size={20} />
          <p class="text-sm text-main font-medium">Verification email sent</p>
          <p class="text-xs text-muted">Check your spam folder if you don't see it.</p>
        </div>
      {:else}
        <p class="text-sm text-muted mb-3">Didn't receive an email?</p>
        <form method="POST" action="?/resend" class="w-full" use:enhance={() => {
          resending = true;
          return async ({ result, update }) => {
            resending = false;
            if (result.type === 'success') {
              resent = true;
            } else if (result.type === 'failure') {
              toast.error((result.data?.error as string) || 'Failed to resend. Please try again.');
            }
            await update({ reset: false });
          };
        }}>
          <Button type="submit" variant="primary" size="lg" class="w-full" disabled={resending}>
            {#if resending}
              <RefreshCw size={14} class="animate-spin" /> Resending…
            {:else}
              Resend verification email
            {/if}
          </Button>
        </form>
      {/if}
    </div>
  </div>

  <!-- Footer link -->
  <p class="text-center text-sm text-muted mt-8">
    <a href="/auth/login" class="inline-flex items-center gap-1.5 font-medium hover:text-primary transition-colors">
      <ArrowLeft size={13} />
      Back to login
    </a>
  </p>
</div>
