<script lang="ts">
  import { page } from '$app/state';
  import { beforeNavigate, afterNavigate } from '$app/navigation';
  import '../app.css';
  import Header from '$lib/components/layout/Header.svelte';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';

  let { children } = $props();

  let isAuthRoute = $derived(['/login', '/signup'].includes(page.url.pathname));

  let isNavigating = $state(false);

  beforeNavigate(() => {
    isNavigating = true;
  });

  afterNavigate(() => {
    isNavigating = false;
  });
</script>

{#if isNavigating}
  <div class="fixed top-0 left-0 w-full h-1 z-[100] bg-void-black overflow-hidden pointer-events-none">
    <div class="h-full bg-electric-lime w-1/3 animate-progress origin-left rounded-r-full shadow-[0_0_10px_rgba(204,255,0,0.8)]"></div>
  </div>
{/if}

{#if !isAuthRoute}
  <div class="flex h-[100dvh] w-screen overflow-hidden bg-void-black text-paper font-inter antialiased">
    <Sidebar />
    <div class="flex flex-col flex-1 min-w-0 overflow-hidden relative">
      <Header />
      <main class="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-10">
        <div class="max-w-6xl mx-auto w-full pb-20">
          {@render children()}
        </div>
      </main>
    </div>
  </div>
{:else}
  <div class="min-h-[100dvh] flex flex-col bg-void-black text-paper font-inter antialiased">
    <main class="flex-grow flex flex-col items-center justify-center p-4">
      {@render children()}
    </main>
  </div>
{/if}

<style>
  @keyframes progress {
    0% { transform: translateX(-100%); width: 30%; }
    50% { width: 40%; }
    100% { transform: translateX(300%); width: 10%; }
  }
  .animate-progress {
    animation: progress 1.5s infinite linear;
  }
</style>
