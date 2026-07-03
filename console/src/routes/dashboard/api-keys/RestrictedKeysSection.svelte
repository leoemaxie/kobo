<script lang="ts">
  import { ShieldAlert, Plus, Lock } from '@lucide/svelte';
</script>

<div>
  <p style="
    font-size: 10px; font-weight: 700; text-transform: uppercase;
    letter-spacing: 0.1em; color: #555; margin: 0 0 10px;
  ">Restricted Keys</p>

  <div style="
    border: 1px dashed #2a2a2a; border-radius: 8px;
    padding: 28px 20px;
    display: flex; align-items: center; justify-content: space-between; gap: 24px;
  ">
    <div style="display: flex; align-items: flex-start; gap: 14px;">
      <ShieldAlert size={16} color="#555" style="margin-top: 2px; flex-shrink: 0;" />
      <div>
        <p style="font-size: 12px; font-weight: 600; color: #888; margin: 0 0 4px;">
          Restricted Keys — not configured
        </p>
        <p style="font-size: 11px; color: #555; margin: 0; line-height: 1.6;">
          Scope keys to specific IP ranges and permission sets for service-to-service use.
        </p>
        <div style="display: flex; align-items: center; gap: 6px; margin-top: 8px;">
          <Lock size={11} color="#444" />
          <span style="font-family: monospace; font-size: 10px; color: #444;">
            scopes: accounts:read, transactions:write
          </span>
        </div>
      </div>
    </div>

    <button style="
      display: flex; align-items: center; gap: 6px; flex-shrink: 0;
      border: 1px solid #2a2a2a; border-radius: 6px;
      background: #111; padding: 6px 12px;
      font-size: 11px; font-weight: 600; color: #888; cursor: pointer;
      white-space: nowrap;
    ">
      <Plus size={12} /> Create restricted key
    </button>
  </div>
</div>
