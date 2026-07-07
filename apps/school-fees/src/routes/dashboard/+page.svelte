<script lang="ts">
  let { data } = $props();
  let students = $derived(data.students);
</script>

<svelte:head>
  <title>Dashboard | Triumph Academy</title>
</svelte:head>

<div class="space-y-8 w-full">
  <div class="flex items-center justify-between border-b border-iron pb-6">
    <div>
      <h1 class="text-3xl font-bold text-pure-white tracking-tight">Dashboard</h1>
      <p class="text-smoke mt-1">Welcome back. Here are your linked students.</p>
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    {#each students as student}
      <a href="/students/{student.id}" class="block group">
        <div class="bg-carbon border border-iron rounded-xl p-6 transition-all hover:border-electric-lime hover:shadow-md hover:bg-graphite/50 h-full flex flex-col">
          <div class="flex items-start justify-between mb-4">
            <div>
              <h2 class="text-xl font-semibold text-pure-white group-hover:text-electric-lime transition-colors">{student.name}</h2>
              <p class="text-sm text-smoke mt-1">{student.class}</p>
            </div>
            <div class="w-10 h-10 rounded-full bg-void-black border border-iron flex items-center justify-center font-bold text-smoke group-hover:border-electric-lime group-hover:text-electric-lime transition-colors">
              {student.name.charAt(0)}
            </div>
          </div>
          
          <div class="mt-auto pt-4 border-t border-iron/50">
            <p class="text-xs text-smoke font-semibold uppercase tracking-widest mb-1">Outstanding Balance</p>
            <p class="text-2xl font-bold text-pure-white font-mono tracking-tight">{student.balance}</p>
          </div>
        </div>
      </a>
    {/each}
  </div>

  {#if students.length === 0}
    <div class="bg-carbon border border-iron border-dashed rounded-xl p-12 text-center">
      <p class="text-smoke">No students linked to your account yet.</p>
      <p class="text-xs text-fog mt-2">Contact the school administrator to link your children.</p>
    </div>
  {/if}
</div>
