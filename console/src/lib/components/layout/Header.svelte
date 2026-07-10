<script lang="ts">
  import { Search, Bell, Slash, Moon, Sun, Menu } from '@lucide/svelte';
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  import { useConsoleState } from '$lib/state/console.svelte';

  let { isMobileMenuOpen = $bindable(false) } = $props();

  const consoleState = useConsoleState();
  let workspaceName = $derived(consoleState.user?.integrator?.name || 'Workspace');
  let workspaceInitial = $derived(workspaceName.charAt(0).toUpperCase());

  // Use global state for environment
  let currentEnv = $derived(consoleState.currentEnvironment);
  let currentTheme = $state('dark');

  onMount(() => {
    currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
  });

  function toggleEnv() {
    consoleState.currentEnvironment =
      consoleState.currentEnvironment === 'sandbox' ? 'production' : 'sandbox';
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
    if (p.includes('kyc')) return 'Compliance';
    return 'Dashboard';
  });
</script>

<header
  class="h-18 border-b border-[var(--border-color)] bg-[var(--bg-header)]/80 backdrop-blur-md flex items-center justify-between px-4 sm:px-7 shrink-0 sticky top-0 z-40 w-full"
>
  <!-- Breadcrumb -->
  <div class="flex items-center gap-2 sm:gap-3 text-[13px] text-muted truncate">
    <button
      class="lg:hidden p-1 -ml-2 text-muted hover:text-main"
      onclick={() => (isMobileMenuOpen = !isMobileMenuOpen)}
      aria-label="Toggle menu"
    >
      <Menu size={20} />
    </button>
    <div
      class="hidden sm:flex h-[22px] w-[22px] rounded-md bg-[var(--bg-active)] border border-[var(--border-color)] items-center justify-center text-[10px] font-extrabold text-main"
    >
      {workspaceInitial}
    </div>
    <span class="hidden sm:inline text-muted px-1.5 py-0.5 rounded-md truncate max-w-[120px]"
      >{workspaceName}</span
    >
    <span class="hidden sm:inline text-subtle text-base font-light mx-0.5">/</span>
    <span
      class="text-main font-semibold px-2 py-0.5 bg-[var(--bg-active)] rounded-md border border-[var(--border-color)] truncate max-w-[150px]"
    >
      {breadcrumb()}
    </span>
  </div>

  <!-- Right side controls -->
  <div class="flex items-center gap-2 sm:gap-3">
    <!-- Search -->
    <div class="relative hidden md:flex items-center">
      <div class="absolute left-2.5 pointer-events-none text-subtle">
        <Search size={13} />
      </div>
      <input
        type="text"
        placeholder="Search..."
        class="bg-[var(--bg-element)] border border-[var(--border-color)] rounded-lg py-2 pl-9 pr-12 text-sm text-main w-64 lg:w-96 outline-none transition-colors focus:border-[var(--border-focus)] placeholder:text-subtle"
      />
      <span
        class="absolute right-2 text-[10px] text-subtle border border-[var(--border-color)] rounded py-[1px] px-[5px] font-mono bg-[var(--bg-active)]"
        >⌘K</span
      >
    </div>

    <div class="hidden md:block h-4 w-[1px] bg-[var(--border-color)]"></div>

    <!-- Env Toggle -->
    <button
      onclick={toggleEnv}
      class="flex items-center gap-1.5 rounded-full border border-[var(--border-color)] bg-[var(--bg-element)] px-2.5 sm:px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider cursor-pointer transition-all hover:bg-[var(--bg-active)]
      {currentEnv === 'sandbox' ? 'text-muted' : 'text-[var(--accent)]'}"
    >
      <span
        class="h-1.5 w-1.5 rounded-full transition-all {currentEnv === 'sandbox'
          ? 'bg-subtle'
          : 'bg-[var(--accent)] drop-shadow-[0_0_4px_var(--accent-glow)]'}"
      ></span>
      <span class="hidden sm:inline">{currentEnv === 'sandbox' ? 'Sandbox' : 'Production'}</span>
      <span class="sm:hidden">{currentEnv === 'sandbox' ? 'SBX' : 'PROD'}</span>
    </button>

    <div class="h-4 w-[1px] bg-[var(--border-color)]"></div>

    <!-- Theme Toggle -->
    <button
      onclick={toggleTheme}
      class="h-9 w-9 rounded-lg border border-transparent bg-transparent flex items-center justify-center text-muted cursor-pointer transition-all hover:bg-[var(--bg-element)] hover:border-[var(--border-color)]"
    >
      {#if currentTheme === 'dark'}
        <Sun size={15} />
      {:else}
        <Moon size={15} />
      {/if}
    </button>

    <!-- Bell -->
    <button
      class="h-9 w-9 rounded-lg border border-transparent bg-transparent flex items-center justify-center text-muted cursor-pointer relative transition-all hover:bg-[var(--bg-element)] hover:border-[var(--border-color)]"
    >
      <Bell size={15} />
      <span
        class="absolute top-1.5 right-1.5 h-1.5 w-1.5 rounded-full bg-[var(--accent)] border-2 border-[var(--bg-app)]"
      ></span>
    </button>

    <!-- Avatar -->
    <div
      class="h-9 w-9 rounded-full bg-[var(--bg-active)] border border-[var(--border-color)] flex items-center justify-center text-xs font-bold text-main cursor-pointer shrink-0"
    >
      {workspaceInitial}
    </div>
  </div>
</header>
