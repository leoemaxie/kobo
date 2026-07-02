<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  let email = $state('');
  let password = $state('');
</script>

<div class="w-full">
  <div class="mb-8 text-center">
    <img src="/logo.png" alt="Kobo Logo" class="h-10 w-auto  mx-auto mb-6" />
    <h1 class="text-2xl font-basier font-semibold text-pure-white">Sign in to Kobo</h1>
    <p class="text-smoke text-sm mt-2">Welcome back to Kobo Console</p>
  </div>

  <Card>
    <form class="space-y-4" method="POST">
      <Input label="Email address" id="email" type="email" bind:value={email} required />
      <Input label="Password" id="password" type="password" bind:value={password} required />
      
      <div class="flex items-center justify-between mt-2">
        <label class="flex items-center text-sm text-smoke">
          <input type="checkbox" class="rounded border-iron bg-void-black text-electric-lime focus:ring-steel h-4 w-4 mr-2" />
          Remember me
        </label>
        <a href="/reset-password" class="text-sm font-medium text-electric-lime hover:text-lime-glow transition-colors">Forgot password?</a>
      </div>

      <div class="pt-2">
        <Button variant="primary" class="w-full" type="submit">Sign In</Button>
      </div>
    </form>
  </Card>

  <p class="text-center text-sm text-smoke mt-6">
    Don't have an account? 
    <a href="/signup" class="font-medium text-electric-lime hover:text-lime-glow transition-colors">Sign up</a>
  </p>
</div>
