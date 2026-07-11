<script lang="ts">
  import { enhance } from '$app/forms';

  let { form } = $props();
  let isSubmitting = $state(false);
</script>

<div class="bg-carbon border border-iron/50 rounded-xl p-6 shadow-sm">
  <h2 class="text-lg font-semibold text-pure-white mb-2">Register Student</h2>
  <p class="text-xs text-smoke mb-6">Creates a Kobo Identity and provisions a virtual account.</p>

  {#if form?.success}
    <div
      class="bg-dark-olive/20 border border-electric-lime/50 text-electric-lime text-sm p-4 rounded-lg mb-6 shadow-sm transition-all animate-in fade-in slide-in-from-top-2"
    >
      Student registered successfully! Kobo Identity created.
    </div>
  {/if}
  {#if form?.error}
    <div
      class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg mb-6 shadow-sm transition-all animate-in fade-in slide-in-from-top-2"
    >
      {form.error}
    </div>
  {/if}

  <form
    method="POST"
    action="?/register"
    class="space-y-5"
    use:enhance={() => {
      isSubmitting = true;
      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}
  >
    <div class="space-y-1.5">
      <label
        for="studentId"
        class="block text-xs font-semibold text-smoke uppercase tracking-widest">Student ID</label
      >
      <input
        id="studentId"
        name="studentId"
        type="text"
        required
        placeholder="e.g. STU-12345"
        class="block w-full rounded-lg border border-iron bg-void-black px-4 py-2.5 text-sm text-paper placeholder-fog focus:border-electric-lime/50 focus:outline-none focus:ring-2 focus:ring-electric-lime/20 transition-all duration-300"
      />
    </div>

    <div class="space-y-1.5">
      <label
        for="studentName"
        class="block text-xs font-semibold text-smoke uppercase tracking-widest">Student Name</label
      >
      <input
        id="studentName"
        name="name"
        type="text"
        required
        placeholder="e.g. John Smith"
        class="block w-full rounded-lg border border-iron bg-void-black px-4 py-2.5 text-sm text-paper placeholder-fog focus:border-electric-lime/50 focus:outline-none focus:ring-2 focus:ring-electric-lime/20 transition-all duration-300"
      />
    </div>

    <div class="space-y-1.5">
      <label
        for="className"
        class="block text-xs font-semibold text-smoke uppercase tracking-widest"
        >Class / Grade</label
      >
      <input
        id="className"
        name="className"
        type="text"
        required
        placeholder="e.g. Grade 10"
        class="block w-full rounded-lg border border-iron bg-void-black px-4 py-2.5 text-sm text-paper placeholder-fog focus:border-electric-lime/50 focus:outline-none focus:ring-2 focus:ring-electric-lime/20 transition-all duration-300"
      />
    </div>

    <div class="pt-2">
      <button
        type="submit"
        disabled={isSubmitting}
        class="w-full rounded-lg bg-electric-lime text-void-black px-4 py-2.5 text-sm font-bold shadow-md shadow-electric-lime/10 hover:shadow-electric-lime/20 hover:-translate-y-0.5 active:translate-y-0 active:scale-[0.98] transition-all duration-300 disabled:opacity-50 disabled:hover:translate-y-0 disabled:hover:shadow-none disabled:active:scale-100"
      >
        {isSubmitting ? 'Creating Identity...' : 'Register Student'}
      </button>
    </div>
  </form>
</div>
