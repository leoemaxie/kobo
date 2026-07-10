<script lang="ts">
  import { ShieldAlert, Plus } from '@lucide/svelte';
</script>

<div
  class="
  rounded-[12px] border border-dashed border-border bg-sidebar/30 p-10
  flex flex-col items-center text-center
"
>
  <div
    class="
    h-14 w-14 rounded-full border border-border bg-background flex items-center
    justify-center mb-5 shadow-sm
  "
  >
    <ShieldAlert size={24} class="text-muted" />
  </div>
  <h3 class="text-base font-semibold text-main mb-2">Restricted Keys</h3>
  <p class="text-sm text-muted max-w-md leading-relaxed mb-6">
    Create API keys with limited permissions and restricted IP addresses to enhance security for
    specific integrations.
  </p>
  <button
    class="
    flex items-center gap-2 rounded-[8px] border border-border bg-element px-5 py-2.5
    text-sm font-semibold text-main hover:bg-element-hover hover:text-main
    transition-all shadow-sm
  "
  >
    <Plus size={14} />
    Create restricted key
  </button>
</div>
