<script lang="ts">
  import { Copy, Check, ChevronRight } from '@lucide/svelte';

  interface Parent {
    id: string;
    name: string;
  }

  interface Student {
    id: string;
    name: string;
    class: string;
    linkedParents?: Parent[];
    virtualAccountNo?: string | null;
    accountName?: string | null;
    date: string;
  }

  let { students }: { students: Student[] } = $props();

  let copiedId = $state<string | null>(null);

  function copyToClipboard(text: string, id: string) {
    navigator.clipboard.writeText(text);
    copiedId = id;
    setTimeout(() => {
      copiedId = null;
    }, 2000);
  }

  function getInitials(name: string) {
    return name.substring(0, 2).toUpperCase();
  }
</script>

<div class="space-y-4">
  <h2 class="text-lg font-semibold text-pure-white mb-2">Registered Students</h2>
  
  <div class="grid grid-cols-1 gap-4">
    {#each students as student}
      <a href="/admin/students/{student.id}" class="block group">
        <div class="bg-carbon border border-iron/50 group-hover:border-electric-lime/30 rounded-xl p-5 transition-all duration-300 group-hover:-translate-y-0.5 shadow-sm group-hover:shadow-electric-lime/5 flex flex-col md:flex-row md:items-center justify-between gap-6">
          
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 shrink-0 rounded-full bg-void-black border border-iron flex items-center justify-center text-sm font-bold text-electric-lime tracking-wider group-hover:border-electric-lime/50 transition-colors">
              {getInitials(student.name)}
            </div>
            <div>
              <div class="text-sm font-medium text-pure-white group-hover:text-electric-lime transition-colors">{student.name}</div>
              <div class="text-xs text-smoke mt-1">Class: {student.class} &bull; Reg: {student.date}</div>
            </div>
          </div>

          <div class="flex flex-col md:flex-row md:items-center gap-4 md:gap-8 flex-1 justify-end">
            
            <div class="hidden lg:block text-right">
              <div class="text-[10px] font-bold uppercase tracking-widest text-smoke mb-1">Parents</div>
              <div class="text-xs text-paper">
                {#if student.linkedParents && student.linkedParents.length > 0}
                  {student.linkedParents.map(p => p.name).join(', ')}
                {:else}
                  <span class="opacity-50 italic">None</span>
                {/if}
              </div>
            </div>

            <div class="text-left md:text-right">
              <div class="text-[10px] font-bold uppercase tracking-widest text-smoke mb-1">Virtual Account</div>
              {#if student.virtualAccountNo}
                <div class="flex items-center md:justify-end gap-2">
                  <span class="text-xs font-medium text-paper">{student.virtualAccountNo}</span>
                  <button
                    type="button"
                    onclick={(e) => { e.preventDefault(); copyToClipboard(student.virtualAccountNo!, student.id); }}
                    class="text-smoke hover:text-electric-lime transition-colors focus:outline-none"
                    title="Copy account number"
                  >
                    {#if copiedId === student.id}
                      <Check size={14} class="text-electric-lime" />
                    {:else}
                      <Copy size={14} />
                    {/if}
                  </button>
                </div>
              {:else}
                <span class="text-xs opacity-50">Not assigned</span>
              {/if}
            </div>

            <div class="shrink-0 text-smoke group-hover:text-electric-lime transition-colors hidden md:block">
              <ChevronRight size={20} />
            </div>

          </div>
        </div>
      </a>
    {/each}

    {#if students.length === 0}
      <div class="bg-carbon border border-iron rounded-xl p-8 text-center">
        <p class="text-sm text-smoke italic">No students registered yet.</p>
      </div>
    {/if}
  </div>
</div>
