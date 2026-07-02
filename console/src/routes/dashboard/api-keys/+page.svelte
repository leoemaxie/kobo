<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Key, Eye, Copy, RefreshCw, Plus } from 'lucide-svelte';

  let keyRevealed = $state(false);
  
  function toggleReveal() {
    keyRevealed = !keyRevealed;
  }
</script>

<div class="space-y-8">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-inter font-bold text-pure-white mb-2 tracking-tight">API Keys</h1>
      <p class="text-bone">Manage your sandbox and production API credentials.</p>
    </div>
  </div>

  <div class="bg-carbon border-l-4 border-electric-lime p-4 rounded-r-radius-md mb-8">
    <div class="flex">
      <div class="flex-shrink-0">
        <Key class="h-5 w-5 text-electric-lime" />
      </div>
      <div class="ml-3">
        <p class="text-sm text-paper">
          You are currently viewing <strong class="text-pure-white font-semibold">Sandbox</strong> keys. 
          Toggle the environment switcher in the navbar to view production keys.
        </p>
      </div>
    </div>
  </div>

  <Card class="p-0 overflow-hidden">
    <div class="p-6 border-b border-iron flex items-center justify-between">
      <div>
        <h2 class="text-lg font-basier font-semibold text-pure-white">Standard Keys</h2>
        <p class="text-sm text-smoke mt-1">These keys will allow you to authenticate API requests.</p>
      </div>
      <Button variant="ghost" class="gap-2 text-sm px-4 py-2">
        <Plus size={16} /> Create secret key
      </Button>
    </div>
    
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-iron">
        <thead class="bg-carbon">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Name</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Key ID</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Secret Key</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium tracking-wider text-smoke uppercase">Created</th>
            <th scope="col" class="relative px-6 py-3"><span class="sr-only">Actions</span></th>
          </tr>
        </thead>
        <tbody class="bg-void-black divide-y divide-iron">
          <tr class="hover:bg-graphite/50 transition-colors">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-paper">Default Sandbox Key</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-smoke font-mono">kobo_sandbox_a1b2c3d4</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-smoke font-mono flex items-center gap-2">
              {#if keyRevealed}
                <span class="text-paper">sk_test_51Nx...8z9Qw</span>
              {:else}
                <span class="text-smoke">••••••••••••••••••••••••</span>
              {/if}
              <button onclick={toggleReveal} class="text-smoke hover:text-electric-lime transition-colors ml-2" aria-label="Toggle reveal key">
                <Eye size={14} />
              </button>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-smoke">Oct 12, 2025</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium flex justify-end gap-3">
              <button class="text-smoke hover:text-electric-lime transition-colors flex items-center gap-1" title="Copy key">
                <Copy size={16} /> <span class="sr-only">Copy</span>
              </button>
              <button class="text-smoke hover:text-electric-lime transition-colors flex items-center gap-1" title="Roll key">
                <RefreshCw size={16} /> <span class="sr-only">Roll</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </Card>

  <Card class="bg-graphite/30 border-dashed border-2">
    <div class="text-center py-6">
      <Key class="mx-auto h-12 w-12 text-smoke mb-3 opacity-50" />
      <h3 class="text-lg font-basier font-semibold text-paper">Restricted Keys</h3>
      <p class="mt-2 text-sm text-smoke max-w-lg mx-auto">
        Create API keys with limited permissions and restricted IP addresses to enhance security for specific integrations.
      </p>
      <div class="mt-6">
        <Button variant="ghost">Create restricted key</Button>
      </div>
    </div>
  </Card>
</div>
