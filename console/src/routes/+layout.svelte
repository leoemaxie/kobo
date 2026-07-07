<script lang="ts">
  import { page, navigating } from '$app/state';
  import '../app.css';
  import Header from '$lib/components/layout/Header.svelte';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';
  import PageSkeleton from '$lib/components/layout/PageSkeleton.svelte';
  import ToastProvider from '$lib/components/ui/ToastProvider.svelte';
  import { initConsoleState } from '$lib/state/console.svelte';

  let { data, children } = $props();

  const consoleState = initConsoleState();

  $effect(() => {
    // This deep-updates all the $state runes inside the class whenever the SvelteKit loader data updates
    consoleState.hydrate(data);
  });

  let isAuthRoute = $derived(
    [
      '/auth/login',
      '/auth/signup',
      '/auth/verify-email',
      '/auth/forgot-password',
      '/auth/reset-password',
      '/dashboard/onboarding'
    ].includes(page.url.pathname)
  );
  let isAdminRoute = $derived(page.url.pathname.startsWith('/admin'));
  let isMobileMenuOpen = $state(false);
</script>

<svelte:head>
  <title>Kobo Console</title>
</svelte:head>
{#if !isAuthRoute && !isAdminRoute}
  <div class="flex h-screen w-screen overflow-hidden bg-[var(--bg-app)]">
    <!-- Overlay for mobile sidebar -->
    {#if isMobileMenuOpen}
      <div 
        class="fixed inset-0 bg-black/50 z-40 lg:hidden"
        onclick={() => isMobileMenuOpen = false}
        role="button"
        tabindex="0"
        aria-label="Close menu"
        onkeydown={(e) => e.key === 'Escape' && (isMobileMenuOpen = false)}
      ></div>
    {/if}

    <!-- Sidebar -->
    <div class="
      fixed inset-y-0 left-0 z-50 transform transition-transform duration-200 ease-in-out
      lg:relative lg:translate-x-0
      {isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}
    ">
      <Sidebar onCloseMobile={() => isMobileMenuOpen = false} />
    </div>

    <!-- Main Content -->
    <div class="flex flex-col flex-1 min-w-0 overflow-hidden w-full">
      <Header bind:isMobileMenuOpen={isMobileMenuOpen} />
      <main class="flex-1 overflow-y-auto p-4 sm:p-6 lg:px-12 lg:py-8 pb-16">
        {#if navigating.to}
          <PageSkeleton />
        {:else}
          {@render children()}
        {/if}
      </main>
    </div>
  </div>
{:else if isAdminRoute}
  <div class="h-screen w-screen overflow-hidden bg-[var(--bg-app)] text-main">
    {@render children()}
  </div>
{:else}
  <div class="min-h-screen flex flex-col">
    <main class="flex-grow flex flex-col items-center justify-center p-4">
      <div class="w-full max-w-md">
        {@render children()}
      </div>
    </main>
  </div>
{/if}

<ToastProvider />
