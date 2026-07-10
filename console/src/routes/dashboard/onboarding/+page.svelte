<script lang="ts">
  import { enhance } from '$app/forms';
  import type { ActionData } from './$types';

  let { form }: { form: ActionData } = $props();
  let loading = $state(false);
</script>

<svelte:head>
  <title>Create Workspace — Kobo Console</title>
</svelte:head>

<div class="w-full text-center">
  <div
    class="mx-auto w-full max-w-[420px] bg-[var(--bg-sidebar)] border border-[#222] rounded-xl py-10 px-9 flex flex-col gap-6"
  >
    <div class="flex items-center justify-center">
      <span
        class="text-[28px] font-extrabold text-[var(--accent)] tracking-[-1px] font-['Inter',sans-serif]"
        >k.</span
      >
    </div>

    <div class="text-center">
      <h1 class="text-[20px] font-bold text-[#f5f5f5] mb-2 tracking-[-0.3px]">
        Create your workspace
      </h1>
      <p class="text-[13px] text-[var(--text-subtle)] m-0 leading-[1.6]">
        Give your project a name. This is how it will appear in the Kobo Console and on API
        responses.
      </p>
    </div>

    {#if form?.error}
      <div
        class="bg-red-500/10 border border-red-500/25 text-red-400 py-2.5 px-3.5 rounded-md text-[13px]"
        role="alert"
      >
        {form.error}
      </div>
    {/if}

    <form
      method="POST"
      class="flex flex-col gap-6"
      use:enhance={() => {
        loading = true;
        return async ({ update }) => {
          await update();
          loading = false;
        };
      }}
    >
      <div class="flex flex-col gap-1.5 text-left">
        <label
          for="name"
          class="text-xs font-semibold text-[var(--text-muted)] uppercase tracking-[0.5px]"
          >Workspace name</label
        >
        <input
          id="name"
          type="text"
          name="name"
          placeholder="e.g. Triumph Systems, Acme Corp"
          required
          minlength="2"
          autocomplete="organization"
          class="bg-[#0f0f0f] border border-[#2a2a2a] rounded-md py-2.5 px-3 text-[#f5f5f5] text-sm focus:outline-none transition-colors duration-150 w-full focus:border-[var(--accent)] placeholder:text-[var(--text-muted)]"
        />
        <span class="text-[11px] text-[var(--text-subtle)]"
          >This will be the display name for your integrator.</span
        >
      </div>

      <button
        type="submit"
        class="w-full p-[11px] bg-[var(--accent)] text-[var(--accent-text)] rounded-md text-sm font-bold transition-all duration-150 hover:opacity-90 hover:-translate-y-[1px] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0"
        disabled={loading}
      >
        {loading ? 'Creating…' : 'Create workspace →'}
      </button>
    </form>

    <p
      class="text-xs text-[var(--text-subtle)] text-center m-0 [&>strong]:text-[var(--text-muted)]"
    >
      You can rename your workspace later in <strong>Settings</strong>.
    </p>
  </div>
</div>
