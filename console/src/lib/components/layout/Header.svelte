<script lang="ts">
  import { Search, Bell, Slash, Moon, Sun } from '@lucide/svelte';
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  const consoleState = useConsoleState();
  let workspaceName = $derived(consoleState.user?.integrator?.name || 'Workspace');
  let workspaceInitial = $derived(workspaceName.charAt(0).toUpperCase());

  let currentEnv = $state('sandbox');
  let currentTheme = $state('dark');

  onMount(() => {
    currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
  });

  function toggleEnv() {
    currentEnv = currentEnv === 'sandbox' ? 'production' : 'sandbox';
  }

  function toggleTheme() {
    currentTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', currentTheme);
    localStorage.setItem('theme', currentTheme);
  }

  let breadcrumb = $derived(() => {
    const p = page.url.pathname;
    if (p === '/dashboard') return 'Overview';
    if (p.includes('api-keys')) return 'API Keys';
    if (p.includes('billing')) return 'Billing';
    if (p.includes('webhooks')) return 'Webhooks';
    if (p.includes('team')) return 'Team';
    if (p.includes('settings')) return 'Settings';
    return 'Dashboard';
  });
</script>

<header style="
  height: 64px; border-bottom: 1px solid var(--border-color);
  background: var(--bg-header); backdrop-filter: blur(12px);
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 28px; flex-shrink: 0; position: sticky; top: 0; z-index: 40;
">
  <!-- Breadcrumb -->
  <div style="display: flex; align-items: center; gap: 6px; font-size: 13px; color: var(--text-muted);">
    <div style="
      height: 22px; width: 22px; border-radius: 5px; background: var(--bg-active);
      border: 1px solid var(--border-color); display: flex; align-items: center; justify-content: center;
      font-size: 10px; font-weight: 800; color: var(--text-main);
    ">{workspaceInitial}</div>
    <span style="color: var(--text-muted); padding: 2px 6px; border-radius: 5px;">{workspaceName}</span>
    <span style="color: var(--text-subtle); font-size: 16px; font-weight: 300; margin: 0 2px;">/</span>
    <span style="
      color: var(--text-main); font-weight: 600; padding: 2px 8px; 
      background: var(--bg-active); border-radius: 5px; border: 1px solid var(--border-color);
    ">{breadcrumb()}</span>
  </div>

  <!-- Right side controls -->
  <div style="display: flex; align-items: center; gap: 12px;">
    <!-- Search -->
    <div style="position: relative; display: flex; align-items: center;">
      <div style="position: absolute; left: 10px; pointer-events: none; color: var(--text-subtle);">
        <Search size={13} />
      </div>
      <input
        type="text"
        placeholder="Search..."
        style="
          background: var(--bg-element); border: 1px solid var(--border-color); border-radius: 8px;
          padding: 6px 40px 6px 30px; font-size: 12px; color: var(--text-main);
          width: 180px; outline: none; transition: border-color 0.2s;
        "
        onfocus={(e) => (e.target as HTMLInputElement).style.borderColor = 'var(--border-focus)'}
        onblur={(e) => (e.target as HTMLInputElement).style.borderColor = 'var(--border-color)'}
      />
      <span style="
        position: absolute; right: 8px; font-size: 10px; color: var(--text-subtle);
        border: 1px solid var(--border-color); border-radius: 4px; padding: 1px 5px;
        font-family: monospace; background: var(--bg-active);
      ">⌘K</span>
    </div>

    <div style="height: 18px; width: 1px; background: var(--border-color);"></div>

    <!-- Env Toggle -->
    <button
      onclick={toggleEnv}
      style="
        display: flex; align-items: center; gap: 7px; border-radius: 99px;
        border: 1px solid var(--border-color); background: var(--bg-element); padding: 5px 12px;
        font-size: 10px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.08em; cursor: pointer; transition: all 0.2s;
        color: {currentEnv === 'sandbox' ? 'var(--text-muted)' : 'var(--accent)'};
      "
    >
      <span style="
        height: 7px; width: 7px; border-radius: 50%;
        background: {currentEnv === 'sandbox' ? 'var(--text-subtle)' : 'var(--accent)'};
        box-shadow: {currentEnv === 'production' ? '0 0 8px var(--accent-glow)' : 'none'};
        display: inline-block; transition: all 0.2s;
      "></span>
      {currentEnv === 'sandbox' ? 'Sandbox' : 'Production'}
    </button>

    <div style="height: 18px; width: 1px; background: var(--border-color);"></div>

    <!-- Theme Toggle -->
    <button
      onclick={toggleTheme}
      style="
        height: 32px; width: 32px; border-radius: 8px; border: 1px solid transparent;
        background: transparent; display: flex; align-items: center; justify-content: center;
        color: var(--text-muted); cursor: pointer; transition: all 0.2s;
      "
      onmouseenter={(e) => {
        (e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-element)';
        (e.currentTarget as HTMLButtonElement).style.borderColor = 'var(--border-color)';
      }}
      onmouseleave={(e) => {
        (e.currentTarget as HTMLButtonElement).style.background = 'transparent';
        (e.currentTarget as HTMLButtonElement).style.borderColor = 'transparent';
      }}
    >
      {#if currentTheme === 'dark'}
        <Sun size={15} />
      {:else}
        <Moon size={15} />
      {/if}
    </button>

    <!-- Bell -->
    <button style="
      height: 32px; width: 32px; border-radius: 8px; border: 1px solid transparent;
      background: transparent; display: flex; align-items: center; justify-content: center;
      color: var(--text-muted); cursor: pointer; position: relative; transition: all 0.2s;
    "
    onmouseenter={(e) => {
      (e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-element)';
      (e.currentTarget as HTMLButtonElement).style.borderColor = 'var(--border-color)';
    }}
    onmouseleave={(e) => {
      (e.currentTarget as HTMLButtonElement).style.background = 'transparent';
      (e.currentTarget as HTMLButtonElement).style.borderColor = 'transparent';
    }}
    >
      <Bell size={15} />
      <span style="
        position: absolute; top: 6px; right: 6px; height: 7px; width: 7px;
        border-radius: 50%; background: var(--accent); border: 2px solid var(--bg-app);
      "></span>
    </button>

    <!-- Avatar -->
    <div style="
      height: 32px; width: 32px; border-radius: 50%;
      background: var(--bg-active);
      border: 1px solid var(--border-color); display: flex; align-items: center; justify-content: center;
      font-size: 12px; font-weight: 700; color: var(--text-main); cursor: pointer;
    ">{workspaceInitial}</div>
  </div>
</header>
