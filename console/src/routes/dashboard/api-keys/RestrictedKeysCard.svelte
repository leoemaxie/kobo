<script lang="ts">
  import { ShieldAlert, Plus } from '@lucide/svelte';
</script>

<div class="
  rounded-[12px] border border-dashed border-iron bg-carbon/30 p-10 
  flex flex-col items-center text-center
">
  <div class="
    h-14 w-14 rounded-full border border-iron bg-void-black flex items-center 
    justify-center mb-5 shadow-sm
  ">
    <ShieldAlert size={24} class="text-smoke" />
  </div>
  <h3 class="text-base font-semibold text-paper mb-2">Restricted Keys</h3>
  <p class="text-sm text-smoke max-w-md leading-relaxed mb-6">
    Create API keys with limited permissions and restricted IP addresses 
    to enhance security for specific integrations.
  </p>
  <button class="
    flex items-center gap-2 rounded-[8px] border border-iron bg-carbon px-5 py-2.5 
    text-sm font-semibold text-paper hover:bg-graphite hover:text-pure-white 
    transition-all shadow-sm
  ">
    <Plus size={14} />
    Create restricted key
  </button>
</div>
