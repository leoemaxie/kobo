<script lang="ts">
  import PageHeader from '$lib/components/common/PageHeader.svelte';
  import RegistrationForm from '$lib/components/admin/students/RegistrationForm.svelte';
  import StudentList from '$lib/components/admin/students/StudentList.svelte';

  let { data, form } = $props();
  let students = $derived(data.students);
  let availableParents = $derived(data.availableParents || []);
</script>

<svelte:head>
  <title>Admin Console | Triumph Academy</title>
</svelte:head>

{#snippet adminBadge()}
  <div
    class="bg-carbon border border-electric-lime/30 text-electric-lime text-xs font-semibold px-3 py-1.5 rounded-full tracking-widest uppercase"
  >
    Admin Mode
  </div>
{/snippet}

<div class="space-y-10 w-full">
  <PageHeader
    title="Admin Console"
    subtitle="Register new students and manage accounts."
    rightSection={adminBadge}
  />

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
    <!-- Left Column: Registration Form -->
    <div class="lg:col-span-1 space-y-6">
      <RegistrationForm {form} />
    </div>

    <!-- Right Column: Student List -->
    <div class="lg:col-span-2 space-y-6">
      <StudentList {students} />
    </div>
  </div>
</div>
