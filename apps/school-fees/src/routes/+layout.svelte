<script lang="ts">
  import { page } from '$app/state';
  import '../app.css';
  import Header from '$lib/components/layout/Header.svelte';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';

  let { children } = $props();

  let isAuthRoute = $derived(['/login', '/signup'].includes(page.url.pathname));
</script>

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
