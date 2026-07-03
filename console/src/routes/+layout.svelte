<script lang="ts">
  import { page, navigating } from '$app/state';
  import '../app.css';
  import Header from '$lib/components/layout/Header.svelte';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';
  import PageSkeleton from '$lib/components/layout/PageSkeleton.svelte';

  let { data, children } = $props();

  let isAuthRoute = $derived(
    [
      '/auth/login',
      '/auth/signup',
      '/auth/verify-email',
      '/auth/forgot-password',
      '/auth/reset-password'
    ].includes(page.url.pathname)
  );
</script>

{#if !isAuthRoute}
  <div style="display:flex; height:100vh; width:100vw; overflow:hidden; background:#080808;">
    <Sidebar />
    <div style="display:flex; flex-direction:column; flex:1; min-width:0; overflow:hidden;">
      <Header />
      <main style="flex:1; overflow-y:auto; padding: 2rem 3rem 4rem;">
        {#if navigating}
          <PageSkeleton />
        {:else}
          {@render children()}
        {/if}
      </main>
    </div>
  </div>
{:else}
  <div class="min-h-screen flex flex-col bg-void-black">
    <main class="flex-grow flex flex-col items-center justify-center p-4">
      <div class="w-full max-w-md">
        {@render children()}
      </div>
    </main>
  </div>
{/if}
