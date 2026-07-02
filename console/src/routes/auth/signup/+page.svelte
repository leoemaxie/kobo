<script lang="ts">
  import { Building2, Mail, Lock, Eye, EyeOff, UserPlus } from '@lucide/svelte';

  let companyName = $state('');
  let email = $state('');
  let password = $state('');
  let showPassword = $state(false);
  let loading = $state(false);
</script>

<div class="w-full">
  <!-- Logo + heading -->
  <div class="mb-8 text-center">
    <img src="/logo.png" alt="Kobo" class="h-10 w-auto mx-auto mb-6" />
    <h1 class="text-3xl font-inter font-bold text-pure-white tracking-tight">Create your account</h1>
    <p class="text-smoke text-sm mt-3 leading-relaxed">
      Already have an account?
      <a href="/auth/login" class="font-semibold text-electric-lime hover:text-lime-glow transition-colors ml-1">Sign in</a>
    </p>
  </div>

  <!-- Card -->
  <div class="bg-carbon border border-iron rounded-[10px] px-16 py-6 shadow-sm">
    <form class="space-y-6" method="POST">
      
      <!-- Company field -->
      <div class="space-y-1.5">
        <label for="company" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Company / Integrator Name</label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
            <Building2 size={15} class="text-fog" />
          </div>
          <input
            id="company"
            type="text"
            bind:value={companyName}
            required
            placeholder="Acme Corp"
            class="block w-full rounded-[8px] border border-iron bg-void-black pl-10 pr-4 py-3 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
          />
        </div>
      </div>

      <!-- Email field -->
      <div class="space-y-1.5">
        <label for="email" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Work Email</label>
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

      <!-- Password field -->
      <div class="space-y-1.5">
        <label for="password" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Password</label>
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

      <div class="pt-1">
        <button
          type="submit"
          disabled={loading}
          class="w-full flex items-center justify-center gap-2 rounded-[8px] bg-electric-lime text-void-black px-6 py-2.5 text-sm font-bold tracking-tight shadow-md hover:bg-lime-glow hover:-translate-y-0.5 active:translate-y-0 transition-all focus:outline-none focus:ring-2 focus:ring-electric-lime focus:ring-offset-2 focus:ring-offset-carbon disabled:opacity-60 disabled:cursor-not-allowed"
        >
          <UserPlus size={16} />
          Create Account
        </button>
      </div>

      <p class="text-[11px] text-fog text-center mt-4 leading-relaxed">
        By creating an account, you agree to our Terms of Service and Privacy Policy.
      </p>
    </form>
  </div>
</div>
