<script lang="ts">
  import { User } from '@lucide/svelte';

  let {
    results,
    query,
    onNavigate,
  }: {
    results: { students: any[]; parents: any[] } | null;
    query: string;
    onNavigate: (id: string) => void;
  } = $props();
</script>

<div class="overflow-y-auto flex-1 p-2 custom-scrollbar">
  {#if results}
    {#if results.students.length === 0 && results.parents.length === 0}
      <div class="py-12 text-center text-smoke">
        <p>No results found for "{query}"</p>
      </div>
    {:else}
      {#if results.students.length > 0}
        <div class="px-3 py-2 text-xs font-bold text-smoke uppercase tracking-widest mt-2">
          Students
        </div>
        <ul class="mb-2">
          {#each results.students as student}
            <li>
              <button
                class="w-full text-left px-4 py-3 rounded-xl hover:bg-graphite/40 transition-colors flex items-center justify-between group"
                onclick={() => onNavigate(student.id)}
              >
                <div>
                  <div
                    class="font-medium text-pure-white group-hover:text-electric-lime transition-colors"
                  >
                    {student.name}
                  </div>
                  <div class="text-xs text-smoke mt-0.5">
                    ID: {student.id} &bull; Class: {student.className}
                  </div>
                </div>
                <div
                  class="text-xs text-smoke font-mono opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  ↵ Enter
                </div>
              </button>
            </li>
          {/each}
        </ul>
      {/if}

      {#if results.parents.length > 0}
        <div
          class="px-3 py-2 text-xs font-bold text-smoke uppercase tracking-widest mt-2 border-t border-iron/50 pt-4"
        >
          Parents
        </div>
        <ul>
          {#each results.parents as parent}
            <li>
              <div class="px-4 py-3 rounded-xl flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div
                    class="w-8 h-8 rounded-full bg-void-black border border-iron flex items-center justify-center text-smoke"
                  >
                    <User size={14} />
                  </div>
                  <div>
                    <div class="font-medium text-pure-white">
                      {parent.name}
                    </div>
                    <div class="text-xs text-smoke mt-0.5">
                      {parent.email} &bull; {parent.role}
                    </div>
                  </div>
                </div>
              </div>
            </li>
          {/each}
        </ul>
      {/if}
    {/if}
  {:else if query.length > 0 && query.length < 2}
    <div class="py-8 text-center text-smoke text-sm">Type at least 2 characters to search...</div>
  {:else}
    <div class="py-8 px-6 flex items-center justify-center gap-4 text-smoke text-xs hidden sm:flex">
      <div class="flex items-center gap-1">
        <kbd
          class="bg-void-black border border-iron px-1.5 py-0.5 rounded shadow-sm font-mono font-bold text-[10px]"
          >↑</kbd
        >
        <kbd
          class="bg-void-black border border-iron px-1.5 py-0.5 rounded shadow-sm font-mono font-bold text-[10px]"
          >↓</kbd
        >
        <span class="ml-1">to navigate</span>
      </div>
      <div class="flex items-center gap-1">
        <kbd
          class="bg-void-black border border-iron px-1.5 py-0.5 rounded shadow-sm font-mono font-bold text-[10px]"
          >↵</kbd
        >
        <span class="ml-1">to select</span>
      </div>
      <div class="flex items-center gap-1">
        <kbd
          class="bg-void-black border border-iron px-1.5 py-0.5 rounded shadow-sm font-mono font-bold text-[10px]"
          >ESC</kbd
        >
        <span class="ml-1">to close</span>
      </div>
    </div>
  {/if}
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 10px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background-color: rgba(255, 255, 255, 0.2);
  }
</style>
