<script lang="ts">
  let { children, variant = 'primary', class: className = '', ...rest } = $props();

  const baseStyles = 'inline-flex items-center justify-center rounded-radius-buttons px-6 py-3 font-inter text-base font-semibold transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-void-black';
  
  const variants = {
    primary: 'bg-electric-lime text-void-black shadow-md hover:bg-lime-glow hover:-translate-y-0.5 active:translate-y-0 focus:ring-electric-lime',
    ghost: 'bg-transparent border border-iron text-paper hover:bg-carbon hover:text-pure-white focus:ring-iron',
    pill: 'rounded-radius-pills bg-dark-olive text-electric-lime px-3 py-1.5 text-sm hover:bg-[#5c5f00]',
  };
</script>

<button 
  class="{baseStyles} {variants[variant as keyof typeof variants]} {className}" 
  {...rest}
>
  {@render children()}
</button>
