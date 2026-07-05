<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Mail } from '@lucide/svelte';
</script>

<div class="w-full text-center">
  <div class="mb-8 flex justify-center">
    <div class="h-16 w-16 rounded-full bg-element border-2 border-border flex items-center justify-center">
      <Mail size={24} class="text-primary" />
    </div>
  </div>
  
  <h1 class="text-2xl font-basier font-semibold text-main mb-2">Check your email</h1>
  <p class="text-muted text-sm mt-2 mb-8 max-w-sm mx-auto">
    We sent a verification link to your email address. Please click the link to verify your account and access the console.
  </p>

  <Card>
    <div class="space-y-4">
      <p class="text-sm text-main">Didn't receive an email?</p>
      <Button variant="ghost" class="w-full">Resend Verification Email</Button>
    </div>
  </Card>

  <div class="mt-8">
    <a href="/login" class="text-sm font-medium text-muted hover:text-primary transition-colors">
      ← Back to login
    </a>
  </div>
</div>
