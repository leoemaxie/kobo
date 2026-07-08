<script lang="ts">
  import { page } from '$app/state';
  import { PUBLIC_DOCS_URL } from '$env/static/public';
  import {
    LayoutDashboard, Key, CreditCard, Settings, Users,
    Webhook, BookOpen, LifeBuoy, ChevronDown, ShieldCheck
  } from '@lucide/svelte';
  import { enhance } from '$app/forms';
  import { toast } from '$lib/state/toast.svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  let { onCloseMobile } = $props<{ onCloseMobile?: () => void }>();

  const navItems = [
    { name: 'Overview', path: '/dashboard', icon: LayoutDashboard },
    { name: 'API Keys', path: '/dashboard/api-keys', icon: Key },
    { name: 'Webhooks', path: '/dashboard/webhooks', icon: Webhook },
    { name: 'Billing', path: '/dashboard/billing', icon: CreditCard },
  ];

  const bottomItems = [
    { name: 'Team', path: '/dashboard/team', icon: Users },
    { name: 'Compliance', path: '/dashboard/kyc', icon: ShieldCheck },
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

<aside class="w-64 min-w-[256px] bg-[var(--bg-sidebar)] border-r border-[var(--border-color)] flex flex-col h-full overflow-hidden transition-colors duration-200">
  <!-- Org Selector -->
  <div class="h-18 flex items-center justify-between px-4 border-b border-[var(--border-color)] shrink-0">
    <button class="flex items-center gap-2.5 w-full bg-transparent border border-transparent rounded-lg py-2 px-2.5 cursor-pointer text-left transition-colors hover:bg-[var(--bg-element)] flex-1 overflow-hidden">
      <div class="h-9 w-9 rounded-md bg-[var(--accent)] flex items-center justify-center font-black text-sm text-[var(--accent-text)] shrink-0">
        {workspaceInitial}
      </div>
      <div class="flex-1 overflow-hidden">
        <p class="text-[13px] font-semibold text-main whitespace-nowrap overflow-hidden text-ellipsis m-0">{workspaceName}</p>
        <p class="text-[10px] font-medium text-muted uppercase tracking-[0.08em] mt-0.5">{workspacePlan}</p>
      </div>
      <ChevronDown size={13} color="var(--text-subtle)" />
    </button>
  </div>

  <!-- Main Nav -->
  <nav class="flex-1 overflow-y-auto px-2.5 pt-5 pb-2.5">
    <p class="text-[10px] font-bold uppercase tracking-[0.1em] text-subtle px-2.5 pb-2">Workspace</p>
    {#each navItems as item}
      <a
        href={item.path}
        class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline mb-0.5 text-[13px] font-medium transition-all hover:bg-[var(--bg-element)]
          {isActive(item.path)
            ? 'bg-[var(--bg-active)] text-main border border-[var(--border-color)]'
            : 'bg-transparent text-muted border border-transparent'}"
        onclick={() => onCloseMobile?.()}
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
  <div class="p-2.5 border-t border-[var(--border-color)] shrink-0">
    <p class="text-[10px] font-bold uppercase tracking-[0.1em] text-subtle px-2.5 pt-1.5 pb-2">Account</p>
    {#each bottomItems as item}
      <a
        href={item.path}
        class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline mb-0.5 text-[13px] font-medium transition-all hover:bg-[var(--bg-element)]
          {isActive(item.path)
            ? 'bg-[var(--bg-active)] text-main border border-[var(--border-color)]'
            : 'bg-transparent text-muted border border-transparent'}"
        onclick={() => onCloseMobile?.()}
      >
        <item.icon size={15} color={isActive(item.path) ? 'var(--accent)' : 'var(--text-muted)'} />
        {item.name}
      </a>
    {/each}
    <div class="h-[1px] bg-[var(--border-color)] mx-1.5 my-2"></div>
    <a href={PUBLIC_DOCS_URL} class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline text-[13px] font-medium text-muted border border-transparent transition-colors hover:bg-[var(--bg-element)]">
      <BookOpen size={15} color="var(--text-subtle)" /> Documentation
    </a>
    <a href="/support" class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline text-[13px] font-medium text-muted border border-transparent transition-colors hover:bg-[var(--bg-element)]" onclick={() => onCloseMobile?.()}>
      <LifeBuoy size={15} color="var(--text-subtle)" /> Support
    </a>
    <div class="h-[1px] bg-[var(--border-color)] mx-1.5 my-2"></div>
    <form onsubmit={async (e) => {
      e.preventDefault();
      try {
        const res = await fetch('/auth/logout', { method: 'POST', redirect: 'follow' });
        if (res.ok || res.redirected) {
          toast.success('Successfully logged out.');
          window.location.href = '/auth/login';
        } else {
          toast.error('Logout failed.');
        }
      } catch (err) {
        toast.error('Logout failed.');
      }
    }}>
      <button type="submit" class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg no-underline text-[13px] w-full font-medium text-[var(--error-color)] border border-transparent bg-transparent cursor-pointer text-left transition-colors hover:bg-[var(--error-bg)]">
        <span class="flex items-center justify-center w-[15px] h-[15px] rounded-full border-[1.5px] border-[var(--error-color)] relative">
          <span class="absolute w-[6px] h-[1.5px] bg-[var(--error-color)] -right-[2px]"></span>
        </span>
        Logout
      </button>
    </form>
  </div>
</aside>
