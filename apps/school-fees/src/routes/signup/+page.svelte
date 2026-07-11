<script lang="ts">
  import SignupForm from '$lib/components/signup/SignupForm.svelte';
  let { form } = $props();
</script>

<svelte:head>
  <title>Sign Up | Triumph Academy</title>
</svelte:head>

<div class="w-full max-w-md space-y-8">
  <div class="text-center">
    <div
      class="w-12 h-12 rounded-lg bg-electric-lime text-void-black flex items-center justify-center text-2xl font-bold mx-auto mb-4 shadow-lg shadow-electric-lime/20"
    >
      K
    </div>
    <h1 class="text-3xl font-bold text-pure-white tracking-tight">Create Account</h1>
    <p class="text-smoke mt-2">Register for the Triumph Academy parent portal</p>
  </div>

  <SignupForm {form} />

  <p class="text-center text-sm text-smoke">
    Already have an account?
    <a href="/login" class="font-medium text-electric-lime hover:text-lime-glow transition-colors"
      >Sign in</a
    >
  </p>
</div>
