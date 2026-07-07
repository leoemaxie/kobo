<script lang="ts">
  import { page } from '$app/state';
  import { Users, ShieldAlert, History } from '@lucide/svelte';

  let { children } = $props();
</script>

<div class="flex h-full w-full">
  <aside class="w-64 min-w-[256px] bg-[var(--bg-sidebar)] border-r border-[var(--border-color)] flex flex-col h-full overflow-hidden transition-colors duration-200">
    <div class="h-18 flex items-center px-4 border-b border-[var(--border-color)] shrink-0 gap-2.5 text-main">
      <div class="h-9 w-9 rounded-md bg-[var(--error-color)]/10 flex items-center justify-center font-black text-sm text-[var(--error-color)] border border-[var(--error-color)]/20 shrink-0">
        <ShieldAlert size={16} />
      </div>
      <div class="flex-1 overflow-hidden">
        <p class="text-[13px] font-semibold text-main whitespace-nowrap overflow-hidden text-ellipsis m-0">Kobo Console</p>
        <p class="text-[10px] font-medium text-[var(--error-color)] uppercase tracking-[0.08em] mt-0.5">Superadmin</p>
      </div>
    </div>

    <nav class="flex-1 overflow-y-auto px-2.5 pt-5 pb-2.5">
      <p class="text-[10px] font-bold uppercase tracking-[0.1em] text-subtle px-2.5 pb-2">Administration</p>
      <a href="/admin/integrators" 
         class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline mb-0.5 text-[13px] font-medium transition-all hover:bg-[var(--bg-element)] {page.url.pathname.startsWith('/admin/integrators') ? 'bg-[var(--bg-active)] text-main border border-[var(--border-color)]' : 'bg-transparent text-muted border border-transparent'}">
        <Users size={15} color={page.url.pathname.startsWith('/admin/integrators') ? 'var(--accent)' : 'var(--text-muted)'} /> 
        Integrators
      </a>
      <a href="/admin/audit-log" 
         class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline mb-0.5 text-[13px] font-medium transition-all hover:bg-[var(--bg-element)] {page.url.pathname.startsWith('/admin/audit-log') ? 'bg-[var(--bg-active)] text-main border border-[var(--border-color)]' : 'bg-transparent text-muted border border-transparent'}">
        <History size={15} color={page.url.pathname.startsWith('/admin/audit-log') ? 'var(--accent)' : 'var(--text-muted)'} /> 
        Audit Log
      </a>
      
    </nav>
  </aside>

  <div class="flex flex-col flex-1 min-w-0 overflow-hidden w-full">
    <header class="h-18 border-b border-[var(--border-color)] bg-[var(--bg-header)]/80 backdrop-blur-md flex items-center justify-between px-4 sm:px-7 shrink-0 sticky top-0 z-40 w-full">
      <div class="flex items-center gap-2 sm:gap-3 text-[13px] text-muted truncate">
        <span class="text-main font-semibold px-2 py-0.5 bg-[var(--bg-active)] rounded-md border border-[var(--border-color)] truncate max-w-[200px]">
          {page.url.pathname.includes('audit-log') ? 'System Audit Log' : 'Integrator Management'}
        </span>
      </div>
    </header>

    <main class="flex-1 overflow-y-auto p-4 sm:p-6 lg:px-12 lg:py-8 pb-16">
      {@render children()}
    </main>
  </div>
</div>
