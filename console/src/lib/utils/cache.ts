/**
 * Applies standard cache-control headers to SvelteKit server responses.
 * By default, caches privately in the browser for 60 seconds.
 */
export const withCache = (
  setHeaders: (headers: Record<string, string>) => void,
  maxAge: number = 60,
  staleWhileRevalidate: number = 120,
) => {
  setHeaders({
    "cache-control": `private, max-age=${maxAge}, stale-while-revalidate=${staleWhileRevalidate}`,
  });
};
