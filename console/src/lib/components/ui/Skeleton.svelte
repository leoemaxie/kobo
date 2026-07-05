<script lang="ts">
  let { class: className = '', style = '' } = $props<{ class?: string; style?: string }>();
</script>

<div class={`skeleton ${className}`} {style}></div>

<style>
  .skeleton {
    background-color: var(--border-color);
    border-radius: 6px;
    animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }
  
  @keyframes pulse {
    0%, 100% {
      opacity: 1;
    }
    50% {
      opacity: 0.5;
    }
  }
</style>
