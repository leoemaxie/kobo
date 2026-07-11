<script lang="ts">
  import { enhance } from '$app/forms';
  import { Users, Link as LinkIcon, Unlink } from '@lucide/svelte';

  interface Parent {
    id: string;
    name: string;
    email: string;
  }

  let { linkedParents, availableParents }: {
    linkedParents: Parent[];
    availableParents: { id: string; name: string }[];
  } = $props();

  let isSubmitting = $state(false);
</script>

<div class="bg-carbon border border-iron rounded-xl p-6 h-full shadow-sm flex flex-col">
  <div class="flex items-center gap-2 mb-4">
    <Users size={18} class="text-smoke" />
    <h3 class="text-sm font-semibold text-smoke uppercase tracking-widest">Linked Parents</h3>
  </div>

  <div class="flex-1 space-y-3 overflow-y-auto mb-4">
    {#if linkedParents.length > 0}
      {#each linkedParents as parent}
        <div class="flex items-center justify-between bg-void-black border border-iron/50 rounded-lg p-3">
          <div>
            <div class="text-sm font-medium text-pure-white">{parent.name}</div>
            <div class="text-xs text-smoke">{parent.email}</div>
          </div>
          <form method="POST" action="?/unlinkParent" use:enhance={() => {
            isSubmitting = true;
            return async ({ update }) => {
              await update();
              isSubmitting = false;
            };
          }}>
            <input type="hidden" name="parentId" value={parent.id} />
            <button type="submit" disabled={isSubmitting} class="text-xs text-danger/80 hover:text-danger flex items-center gap-1 transition-colors" title="Unlink Parent">
              <Unlink size={14} />
            </button>
          </form>
        </div>
      {/each}
    {:else}
      <div class="text-sm text-smoke italic py-2">No parents linked to this student yet.</div>
    {/if}
  </div>

  {#if availableParents.length > 0}
    <div class="mt-auto pt-4 border-t border-iron/50">
      <form method="POST" action="?/linkParent" use:enhance={() => {
        isSubmitting = true;
        return async ({ update }) => {
          await update();
          isSubmitting = false;
        };
      }} class="flex items-center gap-2">
        <select name="parentId" required class="flex-1 bg-void-black border border-iron rounded-lg px-3 py-2 text-sm text-paper focus:outline-none focus:border-electric-lime transition-colors">
          <option value="" disabled selected>Select parent to link...</option>
          {#each availableParents as parent}
            <option value={parent.id}>{parent.name}</option>
          {/each}
        </select>
        <button type="submit" disabled={isSubmitting} class="bg-iron/50 hover:bg-iron border border-iron text-pure-white p-2 rounded-lg transition-colors flex items-center justify-center disabled:opacity-50">
          <LinkIcon size={18} />
        </button>
      </form>
    </div>
  {/if}
</div>
