<script lang="ts">
  import { page } from '$app/state';
  import { Bell, Search, Menu } from '@lucide/svelte';
  import { ui } from '$lib/state.svelte';
  
  let env = $state('Parent Portal');

  // Reactively switch environment badge text based on route
  $effect(() => {
    if (page.url.pathname.startsWith('/admin')) {
      env = 'Admin Mode';
    } else {
      env = 'Parent Portal';
    }
  });
</script>

<header class="h-16 border-b border-iron bg-void-black/80 backdrop-blur-md sticky top-0 z-30 flex items-center justify-between px-4 sm:px-8 flex-shrink-0">
  <div class="flex items-center gap-4">
    <!-- Mobile Menu Toggle -->
    <button 
      class="lg:hidden text-smoke hover:text-pure-white transition-colors" 
      onclick={() => ui.sidebarOpen = true}
      aria-label="Open Menu"
    >
      <Menu size={24} />
    </button>
    
    <!-- Search Bar -->
    <div class="relative hidden sm:block">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <Search size={14} class="text-fog" />
      </div>
      <input 
        type="text" 
        placeholder="Search..." 
        class="block w-48 lg:w-64 rounded-full border border-iron bg-carbon pl-9 pr-4 py-1.5 text-xs text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
      />
    </div>
  </div>

  <div class="flex items-center gap-3 sm:gap-5">
    <div class="hidden sm:block">
      <span class="inline-flex items-center gap-1.5 rounded-full border border-electric-lime/30 bg-electric-lime/10 px-3 py-1 text-[10px] font-bold uppercase tracking-widest text-electric-lime">
        <span class="w-1.5 h-1.5 rounded-full bg-electric-lime shadow-[0_0_8px_rgba(204,255,0,0.8)]"></span>
        {env}
      </span>
    </div>
    
    <button class="relative text-smoke hover:text-pure-white transition-colors p-2 sm:p-0">
      <Bell size={18} />
      <span class="absolute top-1 sm:top-0 right-1 sm:right-0 w-2 h-2 rounded-full bg-electric-lime border border-void-black"></span>
    </button>
    
    <div class="h-6 w-px bg-iron hidden sm:block"></div>
    
    <div class="flex items-center gap-3">
      <div class="w-8 h-8 rounded-full bg-carbon border border-iron flex items-center justify-center text-sm font-bold text-smoke flex-shrink-0">
        {page.data.user?.name?.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase() || 'U'}
      </div>
      <div class="hidden md:block">
        <p class="text-xs font-semibold text-pure-white leading-none">{page.data.user?.name || 'User'}</p>
        <p class="text-[10px] text-smoke mt-1 truncate max-w-[120px]">{page.data.user?.email || ''}</p>
      </div>
    </div>
  </div>
</header>
