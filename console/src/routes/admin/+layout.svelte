<script lang="ts">
  import { page } from '$app/state';
  import { Users, ShieldAlert, History } from '@lucide/svelte';

  let { children } = $props();
</script>

<div class="flex flex-col md:flex-row gap-8">
  <aside class="w-full md:w-64 flex-shrink-0">
    <div class="sticky top-24">
      <h2 class="text-xs font-semibold text-smoke uppercase tracking-widest mb-4 px-3">Superadmin</h2>
      <nav class="space-y-1">
        <a href="/admin/integrators" 
           class="flex items-center gap-3 px-3 py-2 rounded-radius-md text-sm font-medium transition-colors {page.url.pathname.startsWith('/admin/integrators') ? 'bg-carbon text-pure-white border border-iron' : 'text-smoke hover:bg-carbon hover:text-bone'}">
          <Users size={18} /> Integrators
        </a>
        <a href="/admin/audit-log" 
           class="flex items-center gap-3 px-3 py-2 rounded-radius-md text-sm font-medium transition-colors {page.url.pathname.startsWith('/admin/audit-log') ? 'bg-carbon text-pure-white border border-iron' : 'text-smoke hover:bg-carbon hover:text-bone'}">
          <History size={18} /> Audit Log
        </a>
      </nav>
      
      <div class="mt-8 px-3 py-4 border-t border-iron">
        <div class="flex items-center gap-2 text-electric-lime text-sm font-medium">
          <ShieldAlert size={16} /> Highly privileged area
        </div>
        <p class="text-xs text-smoke mt-2">All actions performed here are recorded in the immutable audit log.</p>
      </div>
    </div>
  </aside>

  <div class="flex-grow">
    {@render children()}
  </div>
</div>
