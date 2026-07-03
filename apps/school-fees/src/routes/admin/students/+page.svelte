<script lang="ts">
  import { enhance } from '$app/forms';

  let { data, form } = $props();
  let students = $derived(data.students);

  let isSubmitting = $state(false);
</script>

<div class="space-y-10 w-full">
  <div class="flex items-center justify-between border-b border-iron pb-6">
    <div>
      <h1 class="text-3xl font-bold text-pure-white tracking-tight">Admin Console</h1>
      <p class="text-smoke mt-1">Register new students and manage accounts.</p>
    </div>
    <div class="bg-carbon border border-electric-lime/30 text-electric-lime text-xs font-semibold px-3 py-1.5 rounded-full tracking-widest uppercase">
      Admin Mode
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
    <!-- Left Column: Registration Form -->
    <div class="lg:col-span-1 space-y-6">
      <div class="bg-carbon border border-iron rounded-xl p-6">
        <h2 class="text-lg font-semibold text-pure-white mb-4">Register Student</h2>
        <p class="text-xs text-smoke mb-6">Creates a Kobo Identity and provisions a virtual account.</p>
        
        {#if form?.success}
          <div class="bg-dark-olive/20 border border-electric-lime/50 text-electric-lime text-sm p-4 rounded-lg mb-6 shadow-sm">
            Student registered successfully! Kobo Identity created.
          </div>
        {/if}
        {#if form?.error}
          <div class="bg-danger/10 border border-danger/50 text-danger text-sm p-4 rounded-lg mb-6 shadow-sm">
            {form.error}
          </div>
        {/if}

        <form method="POST" action="?/register" class="space-y-4" use:enhance={() => {
          isSubmitting = true;
          return async ({ update }) => {
            await update();
            isSubmitting = false;
          };
        }}>
          <div class="space-y-1.5">
            <label for="studentName" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Student Name</label>
            <input
              id="studentName"
              name="name"
              type="text"
              required
              placeholder="e.g. John Smith"
              class="block w-full rounded-lg border border-iron bg-void-black px-4 py-2.5 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
            />
          </div>

          <div class="space-y-1.5">
            <label for="className" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Class / Grade</label>
            <input
              id="className"
              name="className"
              type="text"
              required
              placeholder="e.g. Grade 10"
              class="block w-full rounded-lg border border-iron bg-void-black px-4 py-2.5 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
            />
          </div>

          <div class="pt-2">
            <button
              type="submit"
              disabled={isSubmitting}
              class="w-full rounded-lg bg-electric-lime text-void-black px-4 py-2.5 text-sm font-bold shadow-md hover:bg-lime-glow transition-all disabled:opacity-50"
            >
              {isSubmitting ? 'Creating Identity...' : 'Register Student'}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Right Column: Student List -->
    <div class="lg:col-span-2 space-y-6">
      <h2 class="text-lg font-semibold text-pure-white">Registered Students</h2>
      <div class="bg-carbon border border-iron rounded-xl overflow-hidden shadow-sm">
        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead>
              <tr class="border-b border-iron/50">
                <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Name</th>
                <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Class</th>
                <th class="px-6 py-4 text-left text-[10px] font-bold uppercase tracking-widest text-smoke">Registered</th>
                <th class="px-6 py-4 text-right text-[10px] font-bold uppercase tracking-widest text-smoke">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-iron/50">
              {#each students as student}
                <tr class="hover:bg-graphite/30 transition-colors">
                  <td class="px-6 py-4 text-sm font-medium text-pure-white whitespace-nowrap">{student.name}</td>
                  <td class="px-6 py-4 text-sm text-smoke whitespace-nowrap">{student.class}</td>
                  <td class="px-6 py-4 text-sm text-smoke whitespace-nowrap">{student.date}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-right">
                    <form method="POST" action="?/closeAccount" use:enhance>
                      <input type="hidden" name="studentId" value={student.id} />
                      <button type="submit" class="text-xs font-semibold text-danger border border-danger/30 hover:bg-danger/10 px-3 py-1.5 rounded transition-colors">
                        Close Account
                      </button>
                    </form>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</div>
