<script lang="ts">
  import { enhance } from '$app/forms';
</script>

<svelte:head>
  <title>Account Pending | Triumph Academy</title>
</svelte:head>

<div class="flex flex-col items-center justify-center min-h-[60vh] text-center max-w-lg mx-auto space-y-6">
  <div class="w-16 h-16 rounded-full bg-carbon border border-iron flex items-center justify-center text-3xl mb-2">
    ⏳
  </div>
  
  <h1 class="text-3xl font-bold text-pure-white tracking-tight">Account Pending Approval</h1>
  
  <p class="text-smoke leading-relaxed">
    Your administrator account has been created successfully, but it is currently awaiting approval from the superadmin.
  </p>
  
  <div class="bg-carbon border border-iron rounded-xl p-6 w-full text-left mt-4 shadow-sm">
    <h2 class="text-sm font-semibold text-pure-white mb-2">What happens next?</h2>
    <ul class="space-y-3 text-sm text-smoke">
      <li class="flex items-start gap-2">
        <span class="text-electric-lime mt-0.5">•</span>
        An email has been dispatched to the superadmin requesting access.
      </li>
      <li class="flex items-start gap-2">
        <span class="text-electric-lime mt-0.5">•</span>
        Once approved, you will be granted access to the Admin Console.
      </li>
      <li class="flex items-start gap-2">
        <span class="text-electric-lime mt-0.5">•</span>
        You can check back later by refreshing this page or signing in again.
      </li>
    </ul>
  </div>

  <form method="POST" action="/logout" use:enhance class="mt-8 w-full">
    <button type="submit" class="w-full rounded-lg bg-void-black border border-iron text-pure-white px-4 py-3 text-sm font-bold shadow-sm hover:border-smoke hover:bg-carbon transition-all">
      Sign Out
    </button>
  </form>
</div>
