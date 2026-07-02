<script lang="ts">
  import { page } from '$app/state';
  import { Shield, CreditCard, LayoutDashboard, Database, Key } from '@lucide/svelte';

  let currentEnv = $state('sandbox');

  function toggleEnv() {
    currentEnv = currentEnv === 'sandbox' ? 'production' : 'sandbox';
  }
</script>

<nav class="sticky top-0 z-50 w-full border-b border-iron bg-void-black/80 backdrop-blur-md">
  <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
    <div class="flex h-16 justify-between items-center">
      <div class="flex items-center gap-8">
        <div class="flex flex-shrink-0 items-center gap-3">
          <img src="/logo.png" alt="Kobo Logo" class="h-8 w-auto invert opacity-90" />
          <span class="text-xl font-bold tracking-tight text-pure-white select-none">Console</span>
        </div>
        
        <div class="hidden sm:flex sm:space-x-1">
          <a href="/dashboard" 
             class="flex items-center gap-2 px-3 py-2 rounded-radius-md text-sm font-medium transition-colors {page.url.pathname === '/dashboard' ? 'bg-carbon text-pure-white border border-iron' : 'text-smoke hover:bg-carbon hover:text-bone'}">
            <LayoutDashboard size={16} /> Dashboard
          </a>
          <a href="/api-keys" 
             class="flex items-center gap-2 px-3 py-2 rounded-radius-md text-sm font-medium transition-colors {page.url.pathname === '/api-keys' ? 'bg-carbon text-pure-white border border-iron' : 'text-smoke hover:bg-carbon hover:text-bone'}">
            <Key size={16} /> API Keys
          </a>
          <a href="/billing" 
             class="flex items-center gap-2 px-3 py-2 rounded-radius-md text-sm font-medium transition-colors {page.url.pathname === '/billing' ? 'bg-carbon text-pure-white border border-iron' : 'text-smoke hover:bg-carbon hover:text-bone'}">
            <CreditCard size={16} /> Billing
          </a>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <!-- Clerk-like environment switcher -->
        <button 
          onclick={toggleEnv}
          class="flex items-center gap-2 rounded-radius-full border border-iron bg-carbon px-3 py-1 text-xs font-semibold uppercase tracking-widest transition-all hover:bg-graphite"
        >
          {#if currentEnv === 'sandbox'}
            <span class="h-2 w-2 rounded-full bg-smoke shadow-[0_0_8px_rgba(160,160,160,0.5)]"></span>
            <span class="text-smoke">Sandbox</span>
          {:else}
            <span class="h-2 w-2 rounded-full bg-electric-lime shadow-[0_0_8px_rgba(250,255,105,0.5)]"></span>
            <span class="text-electric-lime">Production</span>
          {/if}
        </button>

        <div class="h-8 w-8 rounded-full bg-graphite flex items-center justify-center border border-iron text-sm font-medium text-bone">
          A
        </div>
      </div>
    </div>
  </div>
</nav>
