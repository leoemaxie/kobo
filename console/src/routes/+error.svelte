<script lang="ts">
  import { page } from '$app/state';
  import { Map, AlertTriangle } from '@lucide/svelte';

  let isNotFound = $derived(page.status === 404);
</script>

<div class="flex flex-col items-center justify-center py-20 mt-10">
  <div class="max-w-md w-full text-center">
    <!-- Icon -->
    <div class="mb-6 flex justify-center">
      <div
        class="h-20 w-20 rounded-full bg-[var(--bg-element)] border border-[var(--border-color)] flex items-center justify-center shadow-sm"
      >
        {#if isNotFound}
          <Map size={32} class="text-[var(--accent)]" />
        {:else}
          <AlertTriangle size={32} class="text-[var(--error-color)]" />
        {/if}
      </div>
    </div>

    <!-- Status Code -->
    <h1 class="text-6xl font-basier font-bold text-[var(--text-main)] tracking-tight mb-4">
      {page.status}
    </h1>

    <!-- Message -->
    <h2 class="text-xl font-semibold text-[var(--text-main)] mb-3">
      {isNotFound ? 'Page not found' : 'Something went wrong'}
    </h2>

    <p class="text-[var(--text-subtle)] text-sm mb-10 leading-relaxed">
      {#if isNotFound}
        We couldn't find the page you were looking for. It might have been removed, renamed, or it
        didn't exist in the first place.
      {:else}
        {page.error?.message ||
          'An unexpected error occurred while trying to process your request. Our team has been notified.'}
      {/if}
    </p>

    <!-- Actions -->
    <div class="flex flex-col sm:flex-row gap-4 justify-center">
      <button
        onclick={() => window.history.back()}
        class="flex items-center justify-center rounded-[8px] px-6 py-2.5 font-inter text-sm font-semibold transition-all focus:outline-none border border-[var(--border-color)] bg-[var(--bg-element)] text-[var(--text-muted)] hover:bg-[var(--bg-element-hover)] hover:text-[var(--text-main)]"
      >
        Go Back
      </button>
      <a
        href="/dashboard"
        class="flex items-center justify-center rounded-[8px] px-6 py-2.5 font-inter text-sm font-bold tracking-tight shadow-md transition-all focus:outline-none focus:ring-2 focus:ring-[var(--accent)] focus:ring-offset-2 focus:ring-offset-[var(--bg-app)] bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90 hover:-translate-y-0.5 active:translate-y-0"
      >
        Back to Dashboard
      </a>
    </div>
  </div>
</div>
