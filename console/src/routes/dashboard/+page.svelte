<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Activity, Users, CreditCard, ArrowRight } from 'lucide-svelte';
</script>

<div class="space-y-8">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-inter font-bold text-pure-white mb-2 tracking-tight">Overview</h1>
      <p class="text-bone">Welcome back, Kobo Inc.</p>
    </div>
    <Button variant="primary">View Documentation</Button>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <Card class="flex flex-col gap-2">
      <div class="flex items-center justify-between mb-2">
        <h3 class="text-smoke text-sm font-semibold uppercase tracking-widest">Total Transactions</h3>
        <Activity size={18} class="text-electric-lime" />
      </div>
      <div class="text-4xl font-inter font-bold text-pure-white">1,248</div>
      <div class="text-xs text-smoke mt-2 flex items-center gap-1">
        <span class="text-electric-lime">+12.5%</span> from last week
      </div>
    </Card>

    <Card class="flex flex-col gap-2">
      <div class="flex items-center justify-between mb-2">
        <h3 class="text-smoke text-sm font-semibold uppercase tracking-widest">Active Accounts</h3>
        <Users size={18} class="text-bone" />
      </div>
      <div class="text-4xl font-inter font-bold text-pure-white">342</div>
      <div class="text-xs text-smoke mt-2 flex items-center gap-1">
        <span class="text-electric-lime">+4</span> new this week
      </div>
    </Card>

    <Card class="flex flex-col gap-2">
      <div class="flex items-center justify-between mb-2">
        <h3 class="text-smoke text-sm font-semibold uppercase tracking-widest">Current Balance</h3>
        <CreditCard size={18} class="text-bone" />
      </div>
      <div class="text-4xl font-inter font-bold text-pure-white">₦45,000</div>
      <div class="text-xs text-smoke mt-2">
        Next billing date: <span class="text-bone">Oct 1, 2026</span>
      </div>
    </Card>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
    <div class="lg:col-span-2 space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-basier font-semibold text-pure-white">Recent Activity</h2>
      </div>
      <Card class="p-0 overflow-hidden">
        <div class="divide-y divide-iron">
          {#each [1, 2, 3, 4] as i}
            <div class="p-4 sm:p-6 flex items-center justify-between hover:bg-graphite transition-colors">
              <div class="flex items-center gap-4">
                <div class="h-10 w-10 rounded-full bg-dark-olive flex items-center justify-center border border-iron">
                  <Activity size={16} class="text-electric-lime" />
                </div>
                <div>
                  <p class="text-paper font-medium">API Request: Create Virtual Account</p>
                  <p class="text-sm text-smoke">sandbox.api.kobo.dev/v1/accounts</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-paper font-mono">200 OK</p>
                <p class="text-xs text-smoke">2 mins ago</p>
              </div>
            </div>
          {/each}
        </div>
      </Card>
    </div>

    <div class="space-y-6">
      <h2 class="text-xl font-basier font-semibold text-pure-white">Quick Actions</h2>
      <Card class="space-y-4">
        <button class="w-full flex items-center justify-between p-3 rounded-radius-md border border-iron bg-carbon hover:bg-graphite transition-colors text-left group">
          <div>
            <p class="text-paper font-medium group-hover:text-electric-lime transition-colors">Rotate API Keys</p>
            <p class="text-xs text-smoke">Generate new sandbox credentials</p>
          </div>
          <ArrowRight size={16} class="text-smoke group-hover:text-electric-lime transition-colors" />
        </button>
        <button class="w-full flex items-center justify-between p-3 rounded-radius-md border border-iron bg-carbon hover:bg-graphite transition-colors text-left group">
          <div>
            <p class="text-paper font-medium group-hover:text-electric-lime transition-colors">Invite Team Member</p>
            <p class="text-xs text-smoke">Manage access and roles</p>
          </div>
          <ArrowRight size={16} class="text-smoke group-hover:text-electric-lime transition-colors" />
        </button>
      </Card>
    </div>
  </div>
</div>
