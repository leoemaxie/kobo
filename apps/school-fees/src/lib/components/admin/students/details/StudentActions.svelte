<script lang="ts">
  import { enhance } from '$app/forms';
  import { Settings, BellRing, Trash2, Save } from '@lucide/svelte';

  let { studentId, currentClass }: { studentId: string; currentClass: string } = $props();

  let isSubmitting = $state(false);
  let isEditingClass = $state(false);
  let newClass = $state('');

  $effect(() => {
    if (isEditingClass) {
      newClass = currentClass;
    }
  });
</script>

<div class="bg-carbon border border-iron rounded-xl p-6 shadow-sm flex flex-col">
  <div class="flex items-center gap-2 mb-6">
    <Settings size={18} class="text-smoke" />
    <h3 class="text-sm font-semibold text-smoke uppercase tracking-widest">Management Actions</h3>
  </div>

  <div class="space-y-4">
    <!-- Modify Class -->
    <div class="bg-void-black border border-iron/50 rounded-lg p-4">
      <div class="flex items-center justify-between mb-2">
        <div class="text-xs font-semibold text-pure-white">Class / Grade</div>
        {#if !isEditingClass}
          <button onclick={() => isEditingClass = true} class="text-xs text-electric-lime hover:underline focus:outline-none">Edit</button>
        {/if}
      </div>
      {#if isEditingClass}
        <form method="POST" action="?/modifyClass" use:enhance={() => {
          isSubmitting = true;
          return async ({ update }) => {
            await update();
            isSubmitting = false;
            isEditingClass = false;
          };
        }} class="flex items-center gap-2 mt-2">
          <input type="text" name="className" bind:value={newClass} required class="flex-1 bg-carbon border border-iron rounded-lg px-3 py-1.5 text-sm text-paper focus:outline-none focus:border-electric-lime transition-colors" />
          <button type="submit" disabled={isSubmitting} class="bg-electric-lime/20 text-electric-lime hover:bg-electric-lime/30 p-1.5 rounded-lg transition-colors disabled:opacity-50">
            <Save size={16} />
          </button>
          <button type="button" onclick={() => {isEditingClass = false; newClass = currentClass;}} class="text-smoke hover:text-pure-white p-1.5 rounded-lg transition-colors">
            Cancel
          </button>
        </form>
      {:else}
        <div class="text-sm text-smoke">{currentClass}</div>
      {/if}
    </div>

    <!-- Send Reminder -->
    <form method="POST" action="?/resendReminder" use:enhance={() => {
      isSubmitting = true;
      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}>
      <button type="submit" disabled={isSubmitting} class="w-full flex items-center justify-between bg-void-black border border-iron/50 hover:border-iron hover:bg-graphite/20 rounded-lg p-4 transition-all focus:outline-none disabled:opacity-50 text-left">
        <div>
          <div class="text-sm font-medium text-pure-white">Send Payment Reminder</div>
          <div class="text-xs text-smoke mt-0.5">Notifies linked parents of pending fees</div>
        </div>
        <BellRing size={18} class="text-smoke" />
      </button>
    </form>

    <!-- Close Account -->
    <form method="POST" action="?/closeAccount" use:enhance={() => {
      isSubmitting = true;
      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}>
      <button type="submit" disabled={isSubmitting} onclick={(e) => { if(!confirm('Are you sure you want to close this account? This action refunds any balance to the source.')) e.preventDefault(); }} class="w-full flex items-center justify-between bg-danger/10 border border-danger/30 hover:bg-danger/20 hover:border-danger/50 rounded-lg p-4 transition-all focus:outline-none disabled:opacity-50 text-left mt-2">
        <div>
          <div class="text-sm font-medium text-danger">Close Account</div>
          <div class="text-xs text-danger/70 mt-0.5">Revokes virtual account and refunds balance</div>
        </div>
        <Trash2 size={18} class="text-danger" />
      </button>
    </form>
  </div>
</div>
