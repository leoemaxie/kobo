<script lang="ts">
  import { page } from '$app/state';
  import '../app.css';
  import Navbar from '$lib/components/layout/Navbar.svelte';
  
  let { children } = $props();

  let isAuthRoute = $derived(['/login', '/signup', '/verify-email'].includes(page.url.pathname));
</script>

<div class="min-h-screen flex flex-col bg-void-black">
  {#if !isAuthRoute}
    <Navbar />
  {/if}
  <main class="flex-grow flex flex-col {isAuthRoute ? 'items-center justify-center p-4' : 'pt-8 pb-16'}">
    <div class="w-full {isAuthRoute ? 'max-w-md' : 'mx-auto max-w-[1200px] px-4 sm:px-6 lg:px-8'}">
      {@render children()}
    </div>
  </main>
</div>

