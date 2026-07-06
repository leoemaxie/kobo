<script lang="ts">
  import { enhance } from '$app/forms';
  import type { ActionData } from './$types';

  let { form }: { form: ActionData } = $props();
  let loading = $state(false);
</script>

<svelte:head>
  <title>Create Workspace — Kobo Console</title>
</svelte:head>

<div class="onboarding-root">
  <div class="onboarding-card">

    <div class="onboarding-logo">
      <span class="logo-mark">k.</span>
    </div>

    <div class="onboarding-header">
      <h1>Create your workspace</h1>
      <p>Give your project a name. This is how it will appear in the Kobo Console and on API responses.</p>
    </div>

    {#if form?.error}
      <div class="form-error" role="alert">{form.error}</div>
    {/if}

    <form method="POST" use:enhance={() => {
      loading = true;
      return async ({ update }) => {
        await update();
        loading = false;
      };
    }}>
      <div class="field">
        <label for="name">Workspace name</label>
        <input
          id="name"
          type="text"
          name="name"
          placeholder="e.g. Triumph Systems, Acme Corp"
          required
          minlength="2"
          autocomplete="organization"
          autofocus
        />
        <span class="field-hint">This will be the display name for your integrator.</span>
      </div>

      <button type="submit" class="btn-primary" disabled={loading}>
        {loading ? 'Creating…' : 'Create workspace →'}
      </button>
    </form>

    <p class="footer-note">
      You can rename your workspace later in <strong>Settings</strong>.
    </p>
  </div>
</div>

<style>
  .onboarding-root {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #0a0a0a;
    padding: 24px;
  }

  .onboarding-card {
    width: 100%;
    max-width: 420px;
    background: #111;
    border: 1px solid #222;
    border-radius: 12px;
    padding: 40px 36px;
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .onboarding-logo {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .logo-mark {
    font-size: 28px;
    font-weight: 800;
    color: #c6f135;
    letter-spacing: -1px;
    font-family: 'Inter', sans-serif;
  }

  .onboarding-header {
    text-align: center;
  }

  .onboarding-header h1 {
    font-size: 20px;
    font-weight: 700;
    color: #f5f5f5;
    margin: 0 0 8px;
    letter-spacing: -0.3px;
  }

  .onboarding-header p {
    font-size: 13px;
    color: #666;
    margin: 0;
    line-height: 1.6;
  }

  .form-error {
    background: rgba(239, 68, 68, 0.08);
    border: 1px solid rgba(239, 68, 68, 0.25);
    color: #f87171;
    padding: 10px 14px;
    border-radius: 6px;
    font-size: 13px;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .field label {
    font-size: 12px;
    font-weight: 600;
    color: #888;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .field input {
    background: #0f0f0f;
    border: 1px solid #2a2a2a;
    border-radius: 6px;
    padding: 10px 12px;
    color: #f5f5f5;
    font-size: 14px;
    font-family: inherit;
    outline: none;
    transition: border-color 0.15s;
    width: 100%;
    box-sizing: border-box;
  }

  .field input:focus {
    border-color: #c6f135;
  }

  .field input::placeholder {
    color: #444;
  }

  .field-hint {
    font-size: 11px;
    color: #555;
  }

  .btn-primary {
    width: 100%;
    padding: 11px;
    background: #c6f135;
    color: #0a0a0a;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    transition: opacity 0.15s, transform 0.1s;
    font-family: inherit;
  }

  .btn-primary:hover:not(:disabled) {
    opacity: 0.9;
    transform: translateY(-1px);
  }

  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .footer-note {
    font-size: 12px;
    color: #555;
    text-align: center;
    margin: 0;
  }

  .footer-note strong {
    color: #888;
  }
</style>
