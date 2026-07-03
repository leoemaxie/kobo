<script lang="ts">
  import { Search, Bell, Slash } from '@lucide/svelte';
  import { page } from '$app/state';

  let currentEnv = $state('sandbox');

  function toggleEnv() {
    currentEnv = currentEnv === 'sandbox' ? 'production' : 'sandbox';
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
  height: 64px; border-bottom: 1px solid #1e1e1e;
  background: rgba(8,8,8,0.85); backdrop-filter: blur(12px);
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 28px; flex-shrink: 0; position: sticky; top: 0; z-index: 40;
">
  <!-- Breadcrumb -->
  <div style="display: flex; align-items: center; gap: 6px; font-size: 13px; color: #666;">
    <div style="
      height: 22px; width: 22px; border-radius: 5px; background: #1c1c1c;
      border: 1px solid #2a2a2a; display: flex; align-items: center; justify-content: center;
      font-size: 10px; font-weight: 800; color: #F8F8F8;
    ">K</div>
    <span style="color: #666; padding: 2px 6px; border-radius: 5px;">Kobo Inc.</span>
    <span style="color: #333; font-size: 16px; font-weight: 300; margin: 0 2px;">/</span>
    <span style="
      color: #F8F8F8; font-weight: 600; padding: 2px 8px; 
      background: #1c1c1c; border-radius: 5px; border: 1px solid #2a2a2a;
    ">{breadcrumb()}</span>
  </div>

  <!-- Right side controls -->
  <div style="display: flex; align-items: center; gap: 12px;">
    <!-- Search -->
    <div style="position: relative; display: flex; align-items: center;">
      <div style="position: absolute; left: 10px; pointer-events: none; color: #555;">
        <Search size={13} />
      </div>
      <input
        type="text"
        placeholder="Search..."
        style="
          background: #111; border: 1px solid #2a2a2a; border-radius: 8px;
          padding: 6px 40px 6px 30px; font-size: 12px; color: #C8C8C8;
          width: 180px; outline: none;
        "
        onfocus={(e) => (e.target as HTMLInputElement).style.borderColor = '#C0FF00'}
        onblur={(e) => (e.target as HTMLInputElement).style.borderColor = '#2a2a2a'}
      />
      <span style="
        position: absolute; right: 8px; font-size: 10px; color: #555;
        border: 1px solid #2a2a2a; border-radius: 4px; padding: 1px 5px;
        font-family: monospace; background: #0a0a0a;
      ">⌘K</span>
    </div>

    <div style="height: 18px; width: 1px; background: #222;"></div>

    <!-- Env Toggle -->
    <button
      onclick={toggleEnv}
      style="
        display: flex; align-items: center; gap: 7px; border-radius: 99px;
        border: 1px solid #2a2a2a; background: #111; padding: 5px 12px;
        font-size: 10px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.08em; cursor: pointer;
        color: {currentEnv === 'sandbox' ? '#888' : '#C0FF00'};
      "
    >
      <span style="
        height: 7px; width: 7px; border-radius: 50%;
        background: {currentEnv === 'sandbox' ? '#555' : '#C0FF00'};
        box-shadow: {currentEnv === 'production' ? '0 0 8px rgba(192,255,0,0.6)' : 'none'};
        display: inline-block;
      "></span>
      {currentEnv === 'sandbox' ? 'Sandbox' : 'Production'}
    </button>

    <div style="height: 18px; width: 1px; background: #222;"></div>

    <!-- Bell -->
    <button style="
      height: 32px; width: 32px; border-radius: 8px; border: 1px solid transparent;
      background: transparent; display: flex; align-items: center; justify-content: center;
      color: #666; cursor: pointer; position: relative;
    "
    onmouseenter={(e) => {
      (e.currentTarget as HTMLButtonElement).style.background = '#111';
      (e.currentTarget as HTMLButtonElement).style.borderColor = '#2a2a2a';
    }}
    onmouseleave={(e) => {
      (e.currentTarget as HTMLButtonElement).style.background = 'transparent';
      (e.currentTarget as HTMLButtonElement).style.borderColor = 'transparent';
    }}
    >
      <Bell size={15} />
      <span style="
        position: absolute; top: 6px; right: 6px; height: 7px; width: 7px;
        border-radius: 50%; background: #C0FF00; border: 2px solid #080808;
      "></span>
    </button>

    <!-- Avatar -->
    <div style="
      height: 32px; width: 32px; border-radius: 50%;
      background: linear-gradient(135deg, #1c1c1c, #2a2a2a);
      border: 1px solid #333; display: flex; align-items: center; justify-content: center;
      font-size: 12px; font-weight: 700; color: #F8F8F8; cursor: pointer;
    ">A</div>
  </div>
</header>
