<script lang="ts">
  import { Mail, ArrowRight, ArrowLeft } from '@lucide/svelte';

  let email = $state('');
  let loading = $state(false);
  let submitted = $state(false);
</script>

<div class="w-full">
  <!-- Logo + heading -->
  <div class="mb-8 text-center">
    <img src="/logo.png" alt="Kobo" class="h-10 w-auto mx-auto mb-6" />
    <h1 class="text-3xl font-inter font-bold text-pure-white tracking-tight">Reset password</h1>
    <p class="text-smoke text-sm mt-3 leading-relaxed">
      Remembered it?
      <a href="/auth/login" class="font-semibold text-electric-lime hover:text-lime-glow transition-colors ml-1">Sign in</a>
    </p>
  </div>

  <!-- Card -->
  <div class="bg-carbon border border-iron rounded-[10px] px-16 py-8 shadow-sm">
    {#if !submitted}
      <form class="space-y-6" onsubmit={(e) => { e.preventDefault(); submitted = true; }}>
        
        <p class="text-sm text-smoke text-center mb-6 mt-2">
          Enter the email associated with your account and we'll send you a link to reset your password.
        </p>

        <!-- Email field -->
        <div class="space-y-1.5">
          <label for="email" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Email Address</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
              <Mail size={15} class="text-fog" />
            </div>
            <input
              id="email"
              type="email"
              bind:value={email}
              required
              placeholder="you@company.com"
              class="block w-full rounded-[8px] border border-iron bg-void-black pl-10 pr-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
            />
          </div>
        </div>

        <div class="pt-2">
          <button
            type="submit"
            disabled={loading}
            class="w-full flex items-center justify-center gap-2 rounded-[8px] bg-electric-lime text-void-black px-6 py-2.5 text-sm font-bold tracking-tight shadow-md hover:bg-lime-glow hover:-translate-y-0.5 active:translate-y-0 transition-all focus:outline-none focus:ring-2 focus:ring-electric-lime focus:ring-offset-2 focus:ring-offset-carbon disabled:opacity-60 disabled:cursor-not-allowed"
          >
            Send Reset Link
            <ArrowRight size={16} />
          </button>
        </div>
      </form>
    {:else}
      <div class="text-center space-y-6 py-4">
        <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-void-black border border-iron">
          <Mail size={24} class="text-electric-lime" />
        </div>
        <h2 class="text-lg font-bold text-pure-white">Check your email</h2>
        <p class="text-sm text-smoke px-4">
          We sent a password reset link to <span class="text-paper font-medium">{email}</span>
        </p>
        
        <div class="pt-4">
          <button 
            type="button" 
            onclick={() => submitted = false} 
            class="text-sm font-medium text-electric-lime hover:text-lime-glow transition-colors"
          >
            Didn't receive it? Click to resend.
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>
