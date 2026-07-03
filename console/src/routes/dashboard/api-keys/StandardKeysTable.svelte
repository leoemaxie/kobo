<script lang="ts">
  import { Eye, EyeOff, Copy, RefreshCw, Plus } from '@lucide/svelte';

  let keyRevealed = $state(false);

  const keys = [
    {
      name: 'Default Sandbox Key',
      id: 'kobo_sandbox_a1b2c3d4',
      secret: 'sk_test_51Nx4mBLKjE8z9Qw',
      created: '2025-10-12',
      lastUsed: '2h ago',
    },
  ];

  const cols = ['NAME', 'KEY ID', 'SECRET KEY', 'LAST USED', 'CREATED', ''];
</script>

<div>
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    margin-bottom: 10px;
  ">
    <p style="
      font-size: 10px; font-weight: 700; text-transform: uppercase;
      letter-spacing: 0.1em; color: #555; margin: 0;
    ">Standard Keys</p>
    <button style="
      display: flex; align-items: center; gap: 5px;
      border: 1px solid #2a2a2a; border-radius: 6px;
      background: #111; padding: 5px 10px;
      font-size: 11px; font-weight: 600; color: #888; cursor: pointer;
    ">
      <Plus size={12} /> Create secret key
    </button>
  </div>

  <div style="border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden;">
    <!-- Header -->
    <div style="
      display: grid; grid-template-columns: 1fr 1.4fr 1.6fr 90px 100px 72px;
      padding: 9px 16px; background: #0d0d0d; border-bottom: 1px solid #1e1e1e;
    ">
      {#each cols as col}
        <span style="font-size: 11px; font-weight: 700; letter-spacing: 0.1em; color: #444; text-transform: uppercase;">
          {col}
        </span>
      {/each}
    </div>

    {#each keys as k}
      <div style="
        display: grid; grid-template-columns: 1fr 1.4fr 1.6fr 90px 100px 72px;
        padding: 11px 16px; align-items: center; border-bottom: 1px solid #111;
      ">
        <span style="font-size: 14px; font-weight: 500; color: #C8C8C8;">{k.name}</span>

        <code style="font-family: monospace; font-size: 13px; color: #888;">{k.id}</code>

        <div style="display: flex; align-items: center; gap: 8px;">
          <code style="font-family: monospace; font-size: 13px; color: {keyRevealed ? '#C8C8C8' : '#555'};">
            {keyRevealed ? k.secret : '••••••••••••••••••'}
          </code>
          <button
            onclick={() => (keyRevealed = !keyRevealed)}
            style="background: none; border: none; cursor: pointer; color: #555; padding: 0; display: flex;"
            aria-label="Toggle secret"
          >
            {#if keyRevealed}<EyeOff size={13} />{:else}<Eye size={13} />{/if}
          </button>
        </div>

        <span style="font-family: monospace; font-size: 13px; color: #555;">{k.lastUsed}</span>
        <span style="font-family: monospace; font-size: 13px; color: #555;">{k.created}</span>

        <div style="display: flex; align-items: center; gap: 10px; justify-content: flex-end;">
          <button style="background: none; border: none; cursor: pointer; color: #555; padding: 0; display: flex;"
            title="Copy key">
            <Copy size={13} />
          </button>
          <button style="background: none; border: none; cursor: pointer; color: #555; padding: 0; display: flex;"
            title="Roll key">
            <RefreshCw size={13} />
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>
