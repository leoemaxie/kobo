<script lang="ts">
  import { Send, FileText, MessageSquare, ExternalLink } from '@lucide/svelte';
  import { PUBLIC_DOCS_URL } from '$env/static/public';
</script>

<svelte:head>
  <title>Support — Kobo Console</title>
</svelte:head>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Page bar -->
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    padding-bottom: 20px; border-bottom: 1px solid var(--border-subtle);
  ">
    <div>
      <p style="
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: var(--text-subtle); margin: 0 0 6px;
      ">Support & Resources</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: var(--text-subtle);">status:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: var(--accent-transparent); border: 1px solid var(--accent-border);
          border-radius: 4px; padding: 2px 8px; color: var(--accent); letter-spacing: 0.05em;
        ">all_systems_operational</code>
      </div>
    </div>
  </div>

  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px;">
    <!-- Contact Form -->
    <div style="background: var(--bg-element); border: 1px solid var(--border-subtle); border-radius: 8px; overflow: hidden; grid-column: span 2;">
      <div style="padding: 16px 20px; border-bottom: 1px solid var(--border-subtle); background: var(--bg-sidebar);">
        <h3 style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0; display: flex; align-items: center; gap: 8px;">
          <MessageSquare size={16} color="var(--text-muted)" /> Contact Support
        </h3>
        <p style="font-size: 12px; color: var(--text-subtle); margin: 4px 0 0; padding-left: 24px;">We usually respond within 24 hours.</p>
      </div>
      <div style="padding: 20px; display: grid; gap: 16px; max-width: 600px;">
        <div>
          <label for="subject" style="display: block; font-size: 12px; font-weight: 600; color: var(--text-muted); margin-bottom: 6px;">Subject</label>
          <input id="subject" type="text" placeholder="How can we help?" class="support-input" />
        </div>
        <div>
          <label for="message" style="display: block; font-size: 12px; font-weight: 600; color: var(--text-muted); margin-bottom: 6px;">Message</label>
          <textarea id="message" rows="4" placeholder="Provide details about your issue..." class="support-input"></textarea>
        </div>
        <div>
          <button style="display: flex; align-items: center; justify-content: center; gap: 6px; border: 1px solid var(--accent); border-radius: 7px; background: var(--accent); padding: 6px 12px; font-size: 13px; font-weight: 700; color: var(--accent-text); cursor: pointer; width: fit-content; transition: all 0.15s;">
            <Send size={13} /> Send Message
          </button>
        </div>
      </div>
    </div>

    <!-- Quick Links -->
    <a href={PUBLIC_DOCS_URL} class="support-link">
      <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;">
        <FileText size={20} color="var(--accent)" />
        <ExternalLink size={16} color="var(--text-subtle)" />
      </div>
      <h4 style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0 0 6px;">Documentation</h4>
      <p style="font-size: 13px; color: var(--text-muted); margin: 0;">Explore guides, API references, and integration tutorials.</p>
    </a>

    <a href="https://status.kobo.inc" target="_blank" rel="noopener noreferrer" class="support-link">
      <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;">
        <div style="display: flex; align-items: center; justify-content: center; width: 20px; height: 20px; background: var(--accent-transparent); border-radius: 50%;">
          <div style="width: 8px; height: 8px; background: var(--accent); border-radius: 50%;"></div>
        </div>
        <ExternalLink size={16} color="var(--text-subtle)" />
      </div>
      <h4 style="font-size: 14px; font-weight: 600; color: var(--text-main); margin: 0 0 6px;">System Status</h4>
      <p style="font-size: 13px; color: var(--text-muted); margin: 0;">Check real-time system performance and operational status.</p>
    </a>
  </div>
</div>

<style>
  .support-input {
    width: 100%;
    background: var(--bg-sidebar);
    border: 1px solid #2a2a2a;
    border-radius: 6px;
    padding: 8px 12px;
    font-size: 13px;
    color: var(--text-main);
    outline: none;
    resize: vertical;
    transition: border-color 0.2s;
  }
  .support-input:focus {
    border-color: var(--accent);
  }

  .support-link {
    display: block;
    background: var(--bg-element);
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    padding: 20px;
    text-decoration: none;
    transition: all 0.2s;
  }
  .support-link:hover {
    border-color: var(--text-muted);
    background: var(--bg-sidebar);
  }
</style>
