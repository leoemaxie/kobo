<script lang="ts">
  import { page } from '$app/state';
  import { PUBLIC_DOCS_URL } from '$env/static/public';
  import {
    LayoutDashboard, Key, CreditCard, Settings, Users,
    Webhook, BookOpen, LifeBuoy, ChevronDown
  } from '@lucide/svelte';

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
</script>

<aside style="
  width: 240px;
  min-width: 240px;
  background: #0a0a0a;
  border-right: 1px solid #2a2a2a;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
">
  <!-- Org Selector -->
  <div style="
    height: 64px; display: flex; align-items: center; 
    padding: 0 16px; border-bottom: 1px solid #1e1e1e;
  ">
    <button style="
      display: flex; align-items: center; gap: 10px; width: 100%;
      background: transparent; border: 1px solid transparent; border-radius: 8px;
      padding: 8px 10px; cursor: pointer; text-align: left;
    "
    onmouseenter={(e) => (e.currentTarget as HTMLButtonElement).style.background = '#151515'}
    onmouseleave={(e) => (e.currentTarget as HTMLButtonElement).style.background = 'transparent'}
    >
      <div style="
        height: 32px; width: 32px; border-radius: 6px; background: #C0FF00;
        display: flex; align-items: center; justify-content: center;
        font-weight: 900; font-size: 14px; color: #080808; flex-shrink: 0;
      ">K</div>
      <div style="flex: 1; overflow: hidden;">
        <p style="
          font-size: 13px; font-weight: 600; color: #F8F8F8; 
          white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin: 0;
        ">Kobo Inc.</p>
        <p style="
          font-size: 10px; font-weight: 500; color: #888; 
          text-transform: uppercase; letter-spacing: 0.08em; margin: 2px 0 0;
        ">Free Tier</p>
      </div>
      <ChevronDown size={13} color="#666" />
    </button>
  </div>

  <!-- Main Nav -->
  <nav style="flex: 1; overflow-y: auto; padding: 20px 10px 10px;">
    <p style="
      font-size: 10px; font-weight: 700; text-transform: uppercase; 
      letter-spacing: 0.1em; color: #555; padding: 0 10px 8px;
    ">Workspace</p>
    {#each navItems as item}
      <a
        href={item.path}
        style="
          display: flex; align-items: center; gap: 10px; padding: 7px 10px;
          border-radius: 8px; text-decoration: none; margin-bottom: 2px;
          font-size: 13px; font-weight: 500; transition: all 0.15s;
          {isActive(item.path)
            ? 'background: #1c1c1c; color: #F8F8F8; border: 1px solid #2a2a2a;'
            : 'background: transparent; color: #888; border: 1px solid transparent;'}
        "
      >
        <item.icon
          size={15}
          color={isActive(item.path) ? '#C0FF00' : '#666'}
        />
        {item.name}
      </a>
    {/each}
  </nav>

  <!-- Bottom Nav -->
  <div style="padding: 10px; border-top: 1px solid #1e1e1e;">
    <p style="
      font-size: 10px; font-weight: 700; text-transform: uppercase; 
      letter-spacing: 0.1em; color: #555; padding: 6px 10px 8px;
    ">Account</p>
    {#each bottomItems as item}
      <a
        href={item.path}
        style="
          display: flex; align-items: center; gap: 10px; padding: 7px 10px;
          border-radius: 8px; text-decoration: none; margin-bottom: 2px;
          font-size: 13px; font-weight: 500; transition: all 0.15s;
          {isActive(item.path)
            ? 'background: #1c1c1c; color: #F8F8F8; border: 1px solid #2a2a2a;'
            : 'background: transparent; color: #888; border: 1px solid transparent;'}
        "
      >
        <item.icon size={15} color={isActive(item.path) ? '#C0FF00' : '#666'} />
        {item.name}
      </a>
    {/each}
    <div style="height: 1px; background: #1e1e1e; margin: 8px 6px;"></div>
    <a href={PUBLIC_DOCS_URL} style="
      display: flex; align-items: center; gap: 10px; padding: 7px 10px;
      border-radius: 8px; text-decoration: none; font-size: 13px; 
      font-weight: 500; color: #888; border: 1px solid transparent;
    ">
      <BookOpen size={15} color="#555" /> Documentation
    </a>
    <a href="/support" style="
      display: flex; align-items: center; gap: 10px; padding: 7px 10px;
      border-radius: 8px; text-decoration: none; font-size: 13px; 
      font-weight: 500; color: #888; border: 1px solid transparent;
    ">
      <LifeBuoy size={15} color="#555" /> Support
    </a>
  </div>
</aside>
