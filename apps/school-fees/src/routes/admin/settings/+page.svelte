<script lang="ts">
  import PageHeader from '$lib/components/common/PageHeader.svelte';
  import SettingsForm from '$lib/components/admin/settings/SettingsForm.svelte';

  let { form, data } = $props();
</script>

<svelte:head>
  <title>Settings | Triumph Academy</title>
</svelte:head>

{#snippet getKeysBtn()}
  <a 
    href={data.consoleUrl} 
    target="_blank" 
    rel="noopener noreferrer" 
    class="bg-electric-lime/20 border border-electric-lime/50 text-electric-lime hover:bg-electric-lime hover:text-void-black text-xs font-bold px-4 py-2 rounded-lg transition-all shadow-sm flex items-center gap-1 whitespace-nowrap"
  >
    Get Keys in Console ↗
  </a>
{/snippet}

<div class="space-y-10 w-full">
  <PageHeader
    title="Kobo Integration Settings"
    subtitle="Configure API credentials for interacting with Kobo."
    rightSection={getKeysBtn}
  />

  <SettingsForm {form} />
</div>
