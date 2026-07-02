<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  let companyName = $state('');
  let email = $state('');
  let password = $state('');
</script>

<div class="w-full">
  <div class="mb-8 text-center">
    <img src="/logo.png" alt="Kobo Logo" class="h-10 w-auto invert opacity-90 mx-auto mb-6" />
    <h1 class="text-2xl font-basier font-semibold text-pure-white">Create an account</h1>
    <p class="text-smoke text-sm mt-2">Start building with Kobo's sandbox</p>
  </div>

  <Card>
    <form class="space-y-4" method="POST">
      <Input label="Company / Integrator Name" id="company" bind:value={companyName} required />
      <Input label="Work Email" id="email" type="email" bind:value={email} required />
      <Input label="Password" id="password" type="password" bind:value={password} required />
      
      <div class="pt-2">
        <Button variant="primary" class="w-full" type="submit">Create Account</Button>
      </div>
      <p class="text-xs text-fog text-center mt-4">
        By creating an account, you agree to our Terms of Service and Privacy Policy.
      </p>
    </form>
  </Card>

  <p class="text-center text-sm text-smoke mt-6">
    Already have an account? 
    <a href="/login" class="font-medium text-electric-lime hover:text-lime-glow transition-colors">Sign in</a>
  </p>
</div>
