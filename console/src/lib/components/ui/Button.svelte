<script lang="ts">
  /**
   * Button — primary design-system button.
   *
   * Variants:
   *  primary  — accent colour fill (default)
   *  ghost    — transparent with border
   *  danger   — red destructive action
   *  neutral  — sidebar-coloured, for secondary actions
   *  pill     — small rounded pill (e.g. environment tags)
   *
   * Sizes:
   *  sm — compact (13px, py-1.5 px-3)
   *  md — default (13px, py-1.5 px-3) — same as sm for now, used in dashboards
   *  lg — large auth submit button (14px, py-2.5 px-6)
   */
  let {
    children,
    variant = 'primary',
    size = 'md',
    class: className = '',
    ...rest
  }: {
    children: import('svelte').Snippet;
    variant?: 'primary' | 'ghost' | 'danger' | 'neutral' | 'pill';
    size?: 'sm' | 'md' | 'lg';
    class?: string;
    [key: string]: unknown;
  } = $props();

  const base =
    'inline-flex items-center justify-center gap-1.5 font-inter font-bold cursor-pointer transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-background disabled:opacity-60 disabled:cursor-not-allowed';

  const variants: Record<string, string> = {
    primary:
      'bg-primary text-primary-text border border-primary rounded-[7px] shadow-md hover:opacity-90 hover:-translate-y-0.5 active:translate-y-0 focus:ring-primary',
    ghost:
      'bg-transparent border border-border text-main rounded-[7px] hover:bg-element-hover focus:ring-border',
    danger:
      'bg-transparent border border-red-500 text-red-500 rounded-[6px] hover:bg-red-500 hover:text-white focus:ring-red-500',
    neutral:
      'bg-sidebar border border-border text-muted rounded-[6px] hover:text-main focus:ring-border',
    pill: 'bg-element text-primary rounded-full border border-primary-border px-3 py-1 text-xs hover:opacity-90',
  };

  const sizes: Record<string, string> = {
    sm: 'px-2.5 py-1.5 text-[11px]',
    md: 'px-3 py-1.5 text-[13px]',
    lg: 'px-6 py-2.5 text-sm tracking-tight',
  };
</script>

<button class="{base} {variants[variant]} {sizes[size]} {className}" {...rest}>
  {@render children()}
</button>
