<script lang="ts">
  let name = $state('');
  let className = $state('');
  let isSubmitting = $state(false);
  let showSuccess = $state(false);

  // Mock list of registered students
  const students = [
    { id: '1', name: 'Alex Johnson', class: 'Grade 10', date: 'Oct 1, 2025', status: 'active' },
    { id: '2', name: 'Sam Johnson', class: 'Grade 8', date: 'Oct 1, 2025', status: 'active' },
    { id: '3', name: 'Emma Davis', class: 'Grade 12', date: 'Sep 15, 2025', status: 'active' }
  ];

  function handleSubmit(e: Event) {
    e.preventDefault();
    isSubmitting = true;
    
    // Simulate API call to Kobo POST /v1/identities
    setTimeout(() => {
      isSubmitting = false;
      showSuccess = true;
      name = '';
      className = '';
      setTimeout(() => showSuccess = false, 3000);
    }, 800);
  }
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
        
        {#if showSuccess}
          <div class="bg-dark-olive/20 border border-electric-lime/50 text-electric-lime text-sm p-4 rounded-lg mb-6 shadow-sm">
            Student registered successfully! Kobo Identity created.
          </div>
        {/if}

        <form class="space-y-4" onsubmit={handleSubmit}>
          <div class="space-y-1.5">
            <label for="studentName" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Student Name</label>
            <input
              id="studentName"
              type="text"
              bind:value={name}
              required
              placeholder="e.g. John Smith"
              class="block w-full rounded-lg border border-iron bg-void-black px-4 py-2.5 text-sm text-paper placeholder-fog focus:border-electric-lime focus:outline-none focus:ring-1 focus:ring-electric-lime transition-colors"
            />
          </div>

          <div class="space-y-1.5">
            <label for="className" class="block text-xs font-semibold text-smoke uppercase tracking-widest">Class / Grade</label>
            <input
              id="className"
              type="text"
              bind:value={className}
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
                    <button class="text-xs font-semibold text-danger border border-danger/30 hover:bg-danger/10 px-3 py-1.5 rounded transition-colors">
                      Close Account
                    </button>
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
