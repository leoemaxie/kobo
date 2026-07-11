<script lang="ts">
  import { Search, X, Loader2 } from '@lucide/svelte';
  import { fade } from 'svelte/transition';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import SearchResults from './SearchResults.svelte';

  let { isOpen = $bindable(false) } = $props();

  let query = $state('');
  let results = $state<{ students: any[], parents: any[] } | null>(null);
  let isSearching = $state(false);
  let debounceTimer: ReturnType<typeof setTimeout>;

  let inputElement: HTMLInputElement;

  $effect(() => {
    if (isOpen && inputElement) {
      setTimeout(() => inputElement.focus(), 50);
    } else if (!isOpen) {
      query = '';
      results = null;
    }
  });

  function handleInput() {
    clearTimeout(debounceTimer);
    if (query.trim().length < 2) {
      results = null;
      isSearching = false;
      return;
    }

    isSearching = true;
    debounceTimer = setTimeout(async () => {
      try {
        const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
        if (res.ok) {
          results = await res.json();
        }
      } catch (e) {
        console.error('Search failed', e);
      } finally {
        isSearching = false;
      }
    }, 300);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      isOpen = false;
    }
  }

  function navigateToStudent(id: string) {
    isOpen = false;
    const isAdmin = ['admin', 'superadmin'].includes(page.data.user?.role);
    if (isAdmin) {
      goto(`/admin/students/${id}`);
    } else {
      goto(`/students/${id}`);
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
  <div 
    class="fixed inset-0 z-50 flex items-start justify-center pt-16 sm:pt-24 px-4 bg-void-black/80 backdrop-blur-sm"
    transition:fade={{ duration: 150 }}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="absolute inset-0" onclick={() => isOpen = false}></div>
    
    <div class="relative w-full max-w-2xl bg-carbon border border-iron rounded-2xl shadow-2xl overflow-hidden flex flex-col max-h-[80vh]">
      <!-- Search Input -->
      <div class="flex items-center px-4 py-4 border-b border-iron bg-void-black/50">
        <Search size={20} class="text-electric-lime mr-3" />
        <input 
          bind:this={inputElement}
          bind:value={query}
          oninput={handleInput}
          type="text" 
          placeholder="Search for students (by name/ID) or parents..." 
          class="flex-1 bg-transparent border-none text-pure-white placeholder-smoke focus:ring-0 focus:outline-none text-lg"
        />
        {#if isSearching}
          <Loader2 size={18} class="text-smoke animate-spin ml-3" />
        {/if}
        <button 
          onclick={() => isOpen = false}
          class="ml-3 p-1 rounded-md text-smoke hover:text-pure-white hover:bg-iron/50 transition-colors"
        >
          <X size={18} />
        </button>
      </div>

      <!-- Results -->
      <SearchResults {results} {query} onNavigate={navigateToStudent} />
    </div>
  </div>
{/if}
