<script lang="ts">
  let { children, variant = 'primary', class: className = '', ...rest } = $props();

  const baseStyles = 'inline-flex items-center justify-center rounded-radius-buttons px-6 py-3 font-inter text-base font-semibold transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-background';
  
  const variants = {
    primary: 'bg-primary text-primary-text shadow-md hover:opacity-90 hover:-translate-y-0.5 active:translate-y-0 focus:ring-primary',
    ghost: 'bg-transparent border border-border text-main hover:bg-element hover:text-main focus:ring-border',
    pill: 'rounded-radius-pills bg-element text-primary px-3 py-1.5 text-sm hover:opacity-90',
  };
</script>

<button 
  class="{baseStyles} {variants[variant as keyof typeof variants]} {className}" 
  {...rest}
>
  {@render children()}
</button>
