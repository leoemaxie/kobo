<script lang="ts">
  import { Key, Webhook, Users, ArrowUpRight, CheckCircle2 } from '@lucide/svelte';

  const actions = [
    {
      icon: Key,
      label: 'Create API Key',
      sub: 'Generate new sandbox credentials.',
      href: '/dashboard/api-keys',
    },
    {
      icon: Webhook,
      label: 'Configure Webhooks',
      sub: 'Register an endpoint for events.',
      href: '/dashboard/webhooks',
    },
    {
      icon: Users,
      label: 'Invite Teammate',
      sub: 'Add a collaborator to this project.',
      href: '/dashboard/team',
    },
  ];

  const checklist = [
    { done: true,  label: 'Create API key' },
    { done: true,  label: 'Provision a virtual account' },
    { done: false, label: 'Send your first transaction' },
    { done: false, label: 'Register a webhook endpoint' },
    { done: false, label: 'Go live (KYC verification)' },
  ];

  const pct = Math.round((checklist.filter(c => c.done).length / checklist.length) * 100);
</script>

<div class="space-y-6">
  <!-- Quick Links -->
  <div>
    <p style="font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle); margin-bottom: 8px;">
      Quick Links
    </p>
    <div style="border: 1px solid var(--border-color); border-radius: 8px; overflow: hidden;">
      {#each actions as a, i}
        <a href={a.href} style="
          display: flex; align-items: center; justify-content: space-between;
          padding: 12px 14px; text-decoration: none;
          border-bottom: {i < actions.length - 1 ? '1px solid var(--border-color)' : 'none'};
          transition: background 0.1s;
        "
          onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.background = 'var(--bg-element)'}
          onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.background = 'transparent'}
        >
          <div style="display: flex; align-items: center; gap: 10px;">
            <a.icon size={14} color="var(--text-subtle)" />
            <div>
              <p style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0;">{a.label}</p>
              <p style="font-size: 12px; color: var(--text-subtle); margin: 2px 0 0;">{a.sub}</p>
            </div>
          </div>
          <ArrowUpRight size={13} color="var(--text-muted)" />
        </a>
      {/each}
    </div>
  </div>

  <!-- Setup Checklist -->
  <div>
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
      <p style="font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-subtle); margin: 0;">
        Setup Checklist
      </p>
      <span style="font-family: monospace; font-size: 12px; color: var(--accent);">{pct}%</span>
    </div>
    <div style="
      width: 100%; height: 2px; background: var(--border-color);
      border-radius: 2px; margin-bottom: 12px; overflow: hidden;
    ">
      <div style="
        height: 100%; width: {pct}%; background: var(--accent);
        box-shadow: 0 0 6px rgba(192,255,0,0.5); border-radius: 2px;
      "></div>
    </div>
    <div style="border: 1px solid var(--border-color); border-radius: 8px; overflow: hidden;">
      {#each checklist as item, i}
        <div style="
          display: flex; align-items: center; gap: 10px;
          padding: 10px 14px;
          border-bottom: {i < checklist.length - 1 ? '1px solid var(--border-color)' : 'none'};
        ">
          <CheckCircle2 size={13} color={item.done ? 'var(--accent)' : 'var(--text-muted)'} />
          <span style="
            font-size: 13px;
            color: {item.done ? 'var(--text-subtle)' : 'var(--text-main)'};
            text-decoration: {item.done ? 'line-through' : 'none'};
          ">{item.label}</span>
        </div>
      {/each}
    </div>
  </div>
</div>
