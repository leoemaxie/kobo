<script lang="ts">
  import { Send, FileText, MessageSquare, ExternalLink } from '@lucide/svelte';
  import { PUBLIC_DOCS_URL } from '$env/static/public';
</script>

<div style="display: flex; flex-direction: column; gap: 28px;">
  <!-- Page bar -->
  <div style="
    display: flex; align-items: center; justify-content: space-between;
    padding-bottom: 20px; border-bottom: 1px solid #1e1e1e;
  ">
    <div>
      <p style="
        font-size: 12px; font-weight: 700; text-transform: uppercase;
        letter-spacing: 0.1em; color: #555; margin: 0 0 6px;
      ">Support & Resources</p>
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-family: monospace; font-size: 13px; color: #555;">status:</span>
        <code style="
          font-family: monospace; font-size: 11px;
          background: rgba(192,255,0,0.08); border: 1px solid rgba(192,255,0,0.2);
          border-radius: 4px; padding: 2px 8px; color: #C0FF00; letter-spacing: 0.05em;
        ">all_systems_operational</code>
      </div>
    </div>
  </div>

  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px;">
    <!-- Contact Form -->
    <div style="background: #0a0a0a; border: 1px solid #1e1e1e; border-radius: 8px; overflow: hidden; grid-column: span 2;">
      <div style="padding: 16px 20px; border-bottom: 1px solid #1e1e1e; background: #111;">
        <h3 style="font-size: 14px; font-weight: 600; color: #F8F8F8; margin: 0; display: flex; align-items: center; gap: 8px;">
          <MessageSquare size={16} color="#888" /> Contact Support
        </h3>
        <p style="font-size: 12px; color: #666; margin: 4px 0 0; padding-left: 24px;">We usually respond within 24 hours.</p>
      </div>
      <div style="padding: 20px; display: grid; gap: 16px; max-width: 600px;">
        <div>
          <label for="subject" style="display: block; font-size: 12px; font-weight: 600; color: #888; margin-bottom: 6px;">Subject</label>
          <input id="subject" type="text" placeholder="How can we help?" class="support-input" />
        </div>
        <div>
          <label for="message" style="display: block; font-size: 12px; font-weight: 600; color: #888; margin-bottom: 6px;">Message</label>
          <textarea id="message" rows="4" placeholder="Provide details about your issue..." class="support-input"></textarea>
        </div>
        <div>
          <button style="display: flex; align-items: center; justify-content: center; gap: 6px; border: 1px solid #C0FF00; border-radius: 6px; background: #C0FF00; padding: 8px 16px; font-size: 13px; font-weight: 700; color: #080808; cursor: pointer; width: fit-content;">
            <Send size={13} /> Send Message
          </button>
        </div>
      </div>
    </div>

    <!-- Quick Links -->
    <a href={PUBLIC_DOCS_URL} class="support-link">
      <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;">
        <FileText size={20} color="#C0FF00" />
        <ExternalLink size={16} color="#666" />
      </div>
      <h4 style="font-size: 14px; font-weight: 600; color: #F8F8F8; margin: 0 0 6px;">Documentation</h4>
      <p style="font-size: 13px; color: #888; margin: 0;">Explore guides, API references, and integration tutorials.</p>
    </a>

    <a href="https://status.kobo.inc" target="_blank" rel="noopener noreferrer" class="support-link">
      <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;">
        <div style="display: flex; align-items: center; justify-content: center; width: 20px; height: 20px; background: rgba(192,255,0,0.1); border-radius: 50%;">
          <div style="width: 8px; height: 8px; background: #C0FF00; border-radius: 50%;"></div>
        </div>
        <ExternalLink size={16} color="#666" />
      </div>
      <h4 style="font-size: 14px; font-weight: 600; color: #F8F8F8; margin: 0 0 6px;">System Status</h4>
      <p style="font-size: 13px; color: #888; margin: 0;">Check real-time system performance and operational status.</p>
    </a>
  </div>
</div>

<style>
  .support-input {
    width: 100%;
    background: #111;
    border: 1px solid #2a2a2a;
    border-radius: 6px;
    padding: 8px 12px;
    font-size: 13px;
    color: #F8F8F8;
    outline: none;
    resize: vertical;
    transition: border-color 0.2s;
  }
  .support-input:focus {
    border-color: #C0FF00;
  }

  .support-link {
    display: block;
    background: #0a0a0a;
    border: 1px solid #1e1e1e;
    border-radius: 8px;
    padding: 20px;
    text-decoration: none;
    transition: all 0.2s;
  }
  .support-link:hover {
    border-color: #444;
    background: #111;
  }
</style>
