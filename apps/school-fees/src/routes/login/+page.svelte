<script lang="ts">
  import LoginForm from '$lib/components/login/LoginForm.svelte';
  let { form } = $props();
</script>

<svelte:head>
  <title>Login | Triumph Academy</title>
</svelte:head>

<div class="w-full max-w-md space-y-8">
  <div class="text-center">
    <div class="w-12 h-12 rounded-lg bg-electric-lime text-void-black flex items-center justify-center text-2xl font-bold mx-auto mb-4 shadow-lg shadow-electric-lime/20">K</div>
    <h1 class="text-3xl font-bold text-pure-white tracking-tight">Parent Portal</h1>
    <p class="text-smoke mt-2">Sign in to view your child's account statement</p>
  </div>

  <LoginForm {form} />

  <p class="text-center text-sm text-smoke">
    Don't have an account? 
    <a href="/signup" class="font-medium text-electric-lime hover:text-lime-glow transition-colors">Sign up</a>
  </p>
</div>
