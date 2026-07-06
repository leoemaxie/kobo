<script lang="ts">
  import { page } from '$app/state';
  import { PUBLIC_DOCS_URL } from '$env/static/public';
  import {
    LayoutDashboard, Key, CreditCard, Settings, Users,
    Webhook, BookOpen, LifeBuoy, ChevronDown
  } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const navItems = [
    { name: 'Overview', path: '/dashboard', icon: LayoutDashboard },
    { name: 'API Keys', path: '/dashboard/api-keys', icon: Key },
    { name: 'Webhooks', path: '/dashboard/webhooks', icon: Webhook },
    { name: 'Billing', path: '/dashboard/billing', icon: CreditCard },
  ];

  const bottomItems = [
    { name: 'Team', path: '/dashboard/team', icon: Users },
    { name: 'Settings', path: '/dashboard/settings', icon: Settings },
  ];

  function isActive(path: string): boolean {
    if (path === '/dashboard') return page.url.pathname === '/dashboard';
    return page.url.pathname.startsWith(path);
  }

  const consoleState = useConsoleState();
  let workspaceName = $derived(consoleState.user?.integrator?.name || 'Workspace');
  let workspaceInitial = $derived(workspaceName.charAt(0).toUpperCase());
  let workspacePlan = $derived((consoleState.user?.integrator?.plan || 'Free Tier').replace(/_/g, ' '));
</script>

<aside style="
  width: 240px;
  min-width: 240px;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  transition: background 0.2s, border-color 0.2s;
">
  <!-- Org Selector -->
  <div style="
    height: 64px; display: flex; align-items: center; 
    padding: 0 16px; border-bottom: 1px solid var(--border-color);
  ">
    <button style="
      display: flex; align-items: center; gap: 10px; width: 100%;
      background: transparent; border: 1px solid transparent; border-radius: 8px;
      padding: 8px 10px; cursor: pointer; text-align: left; transition: background 0.2s;
    "
    onmouseenter={(e) => (e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-element)'}
    onmouseleave={(e) => (e.currentTarget as HTMLButtonElement).style.background = 'transparent'}
    >
      <div style="
        height: 32px; width: 32px; border-radius: 6px; background: var(--accent);
        display: flex; align-items: center; justify-content: center;
        font-weight: 900; font-size: 14px; color: var(--accent-text); flex-shrink: 0;
      ">{workspaceInitial}</div>
      <div style="flex: 1; overflow: hidden;">
        <p style="
          font-size: 13px; font-weight: 600; color: var(--text-main); 
          white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin: 0;
        ">{workspaceName}</p>
        <p style="
          font-size: 10px; font-weight: 500; color: var(--text-muted); 
          text-transform: uppercase; letter-spacing: 0.08em; margin: 2px 0 0;
        ">{workspacePlan}</p>
      </div>
      <ChevronDown size={13} color="var(--text-subtle)" />
    </button>
  </div>

  <!-- Main Nav -->
  <nav style="flex: 1; overflow-y: auto; padding: 20px 10px 10px;">
    <p style="
      font-size: 10px; font-weight: 700; text-transform: uppercase; 
      letter-spacing: 0.1em; color: var(--text-subtle); padding: 0 10px 8px;
    ">Workspace</p>
    {#each navItems as item}
      <a
        href={item.path}
        style="
          display: flex; align-items: center; gap: 10px; padding: 7px 10px;
          border-radius: 8px; text-decoration: none; margin-bottom: 2px;
          font-size: 13px; font-weight: 500; transition: all 0.15s;
          {isActive(item.path)
            ? 'background: var(--bg-active); color: var(--text-main); border: 1px solid var(--border-color);'
            : 'background: transparent; color: var(--text-muted); border: 1px solid transparent;'}
        "
        onmouseenter={(e) => {
          if (!isActive(item.path)) {
            (e.currentTarget as HTMLAnchorElement).style.background = 'var(--bg-element)';
          }
        }}
        onmouseleave={(e) => {
          if (!isActive(item.path)) {
            (e.currentTarget as HTMLAnchorElement).style.background = 'transparent';
          }
        }}
      >
        <item.icon
          size={15}
          color={isActive(item.path) ? 'var(--accent)' : 'var(--text-muted)'}
        />
        {item.name}
      </a>
    {/each}
  </nav>

  <!-- Bottom Nav -->
  <div style="padding: 10px; border-top: 1px solid var(--border-color);">
    <p style="
      font-size: 10px; font-weight: 700; text-transform: uppercase; 
      letter-spacing: 0.1em; color: var(--text-subtle); padding: 6px 10px 8px;
    ">Account</p>
    {#each bottomItems as item}
      <a
        href={item.path}
        style="
          display: flex; align-items: center; gap: 10px; padding: 7px 10px;
          border-radius: 8px; text-decoration: none; margin-bottom: 2px;
          font-size: 13px; font-weight: 500; transition: all 0.15s;
          {isActive(item.path)
            ? 'background: var(--bg-active); color: var(--text-main); border: 1px solid var(--border-color);'
            : 'background: transparent; color: var(--text-muted); border: 1px solid transparent;'}
        "
        onmouseenter={(e) => {
          if (!isActive(item.path)) {
            (e.currentTarget as HTMLAnchorElement).style.background = 'var(--bg-element)';
          }
        }}
        onmouseleave={(e) => {
          if (!isActive(item.path)) {
            (e.currentTarget as HTMLAnchorElement).style.background = 'transparent';
          }
        }}
      >
        <item.icon size={15} color={isActive(item.path) ? 'var(--accent)' : 'var(--text-muted)'} />
        {item.name}
      </a>
    {/each}
    <div style="height: 1px; background: var(--border-color); margin: 8px 6px;"></div>
    <a href={PUBLIC_DOCS_URL} style="
      display: flex; align-items: center; gap: 10px; padding: 7px 10px;
      border-radius: 8px; text-decoration: none; font-size: 13px; 
      font-weight: 500; color: var(--text-muted); border: 1px solid transparent; transition: background 0.2s;
    "
    onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.background = 'var(--bg-element)'}
    onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.background = 'transparent'}>
      <BookOpen size={15} color="var(--text-subtle)" /> Documentation
    </a>
    <a href="/support" style="
      display: flex; align-items: center; gap: 10px; padding: 7px 10px;
      border-radius: 8px; text-decoration: none; font-size: 13px; 
      font-weight: 500; color: var(--text-muted); border: 1px solid transparent; transition: background 0.2s;
    "
    onmouseenter={(e) => (e.currentTarget as HTMLAnchorElement).style.background = 'var(--bg-element)'}
    onmouseleave={(e) => (e.currentTarget as HTMLAnchorElement).style.background = 'transparent'}>
      <LifeBuoy size={15} color="var(--text-subtle)" /> Support
    </a>
    <div style="height: 1px; background: var(--border-color); margin: 8px 6px;"></div>
    <form method="POST" action="/auth/logout" use:enhance={() => {
      return async ({ result, update }) => {
        if (result.type === 'failure' || result.type === 'error') {
          toast.error('Logout failed.');
        } else {
          toast.success('Successfully logged out.');
        }
        await update();
      };
    }}>
      <button type="submit" style="
        display: flex; align-items: center; gap: 10px; padding: 7px 10px;
        border-radius: 8px; text-decoration: none; font-size: 13px; width: 100%;
        font-weight: 500; color: var(--error-color); border: 1px solid transparent; background: transparent; cursor: pointer; text-align: left; transition: background 0.2s;
      "
      onmouseenter={(e) => (e.currentTarget as HTMLButtonElement).style.background = 'var(--error-bg)'}
      onmouseleave={(e) => (e.currentTarget as HTMLButtonElement).style.background = 'transparent'}
      >
        <span style="display: flex; align-items: center; justify-content: center; width: 15px; height: 15px; border-radius: 50%; border: 1.5px solid var(--error-color); position: relative;">
          <span style="position: absolute; width: 6px; height: 1.5px; background: var(--error-color); right: -2px;"></span>
        </span>
        Logout
      </button>
    </form>
  </div>
</aside>
