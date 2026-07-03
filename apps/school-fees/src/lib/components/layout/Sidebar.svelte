<script lang="ts">
  import { page } from '$app/state';
  import { LayoutDashboard, Users, LogOut, GraduationCap, X } from '@lucide/svelte';
  import { ui } from '$lib/state.svelte';
  
  let navItems = $derived(
    page.data.user?.isAdmin
      ? [{ label: 'Admin Console', href: '/admin/students', icon: Users }]
      : [{ label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard }]
  );
</script>

<!-- Mobile overlay -->
{#if ui.sidebarOpen}
  <div 
    class="fixed inset-0 bg-black/60 z-40 lg:hidden backdrop-blur-sm transition-opacity"
    onclick={() => ui.sidebarOpen = false}
    onkeydown={(e) => e.key === 'Escape' && (ui.sidebarOpen = false)}
    role="button"
    tabindex="0"
    aria-label="Close Sidebar"
  ></div>
{/if}

<!-- Sidebar -->
<aside class="fixed inset-y-0 left-0 z-50 w-64 transform transition-transform duration-300 ease-in-out lg:static lg:translate-x-0 border-r border-iron bg-carbon flex flex-col h-full flex-shrink-0 {ui.sidebarOpen ? 'translate-x-0' : '-translate-x-full'}">
  <div class="h-16 flex items-center justify-between px-6 border-b border-iron">
    <div class="font-bold text-pure-white text-lg flex items-center gap-2">
      <div class="w-7 h-7 rounded bg-electric-lime text-void-black flex items-center justify-center text-sm font-black tracking-tighter shadow-sm shadow-electric-lime/20">
        <GraduationCap size={16} />
      </div>
      <span class="tracking-tight">Triumph Academy</span>
    </div>
    
    <button 
      class="lg:hidden text-smoke hover:text-pure-white p-1 transition-colors" 
      onclick={() => ui.sidebarOpen = false}
      aria-label="Close Menu"
    >
      <X size={20} />
    </button>
  </div>

  <div class="p-4 flex-1 overflow-y-auto">
    <div class="mb-4 px-3 text-[10px] font-bold text-smoke uppercase tracking-widest">
      Menu
    </div>
    <nav class="space-y-1">
      {#each navItems as item}
        {@const isActive = page.url.pathname.startsWith(item.href) && (item.href !== '/dashboard' || page.url.pathname === '/dashboard' || page.url.pathname.startsWith('/students/'))}
        <a 
          href={item.href}
          onclick={() => ui.sidebarOpen = false}
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all {isActive ? 'bg-graphite/50 text-pure-white shadow-sm shadow-black/20' : 'text-smoke hover:bg-void-black hover:text-paper'}"
        >
          <item.icon size={16} class={isActive ? 'text-electric-lime' : ''} />
          {item.label}
        </a>
      {/each}
    </nav>
  </div>

  <div class="p-4 border-t border-iron mt-auto">
    <a href="/login" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-smoke hover:bg-danger/10 hover:text-danger transition-colors">
      <LogOut size={16} />
      Sign Out
    </a>
  </div>
</aside>
