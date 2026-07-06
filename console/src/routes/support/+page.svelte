<script lang="ts">
  import { Send, FileText, MessageSquare, ExternalLink } from '@lucide/svelte';
  import { PUBLIC_DOCS_URL } from '$env/static/public';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import CodeBadge from '$lib/components/ui/CodeBadge.svelte';
  import CardSection from '$lib/components/ui/CardSection.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Textarea from '$lib/components/ui/Textarea.svelte';
  import Button from '$lib/components/ui/Button.svelte';
</script>

<svelte:head>
  <title>Support — Kobo Console</title>
</svelte:head>

<div class="flex flex-col gap-7">
  <PageHeader title="Support &amp; Resources">
    {#snippet meta()}
      <span class="font-inconsolata text-[13px] text-subtle">status:</span>
      <CodeBadge>all_systems_operational</CodeBadge>
    {/snippet}
  </PageHeader>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Contact Form -->
    <div class="col-span-2">
      <CardSection title="Contact Support" subtitle="We usually respond within 24 hours.">
        {#snippet header()}
          <MessageSquare size={16} class="text-muted" />
        {/snippet}
        <div class="grid gap-4 max-w-[600px]">
          <Input id="subject" label="Subject" placeholder="How can we help?" variant="settings" />
          <Textarea id="message" label="Message" rows={4} placeholder="Provide details about your issue..." variant="default" />
          <Button variant="primary" size="md" class="w-fit">
            <Send size={13} /> Send Message
          </Button>
        </div>
      </CardSection>
    </div>

    <!-- Quick Links -->
    <a href={PUBLIC_DOCS_URL} class="block bg-element border border-border-subtle rounded-lg p-5 no-underline hover:border-muted hover:bg-sidebar transition-all">
      <div class="flex items-center justify-between mb-3">
        <FileText size={20} class="text-primary" />
        <ExternalLink size={16} class="text-subtle" />
      </div>
      <h4 class="text-sm font-semibold text-main mb-1.5">Documentation</h4>
      <p class="text-[13px] text-muted m-0">Explore guides, API references, and integration tutorials.</p>
    </a>

    <a href="https://status.kobo.inc" target="_blank" rel="noopener noreferrer" class="block bg-element border border-border-subtle rounded-lg p-5 no-underline hover:border-muted hover:bg-sidebar transition-all">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center justify-center w-5 h-5 bg-primary-transparent rounded-full">
          <div class="w-2 h-2 bg-primary rounded-full"></div>
        </div>
        <ExternalLink size={16} class="text-subtle" />
      </div>
      <h4 class="text-sm font-semibold text-main mb-1.5">System Status</h4>
      <p class="text-[13px] text-muted m-0">Check real-time system performance and operational status.</p>
    </a>
  </div>
</div>
