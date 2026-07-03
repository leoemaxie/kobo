import type { Handle } from '@sveltejs/kit';
import { validateSession } from '$lib/server/auth/session';

const PUBLIC_ROUTES = ['/login', '/signup'];

export const handle: Handle = async ({ event, resolve }) => {
  const sessionId = event.cookies.get('session');
  const result = sessionId ? await validateSession(sessionId) : null;

  event.locals.user = result?.parents ?? null;
  event.locals.session = result?.sessions ?? null;

  const path = event.url.pathname;
  const isPublic = PUBLIC_ROUTES.some((p) => path.startsWith(p));

  if (!isPublic && !event.locals.user) {
    return Response.redirect(new URL('/login', event.url), 302);
  }

  if (path.startsWith('/admin') && !event.locals.user?.isAdmin) {
    return new Response('Not found', { status: 404 });
  }

  return resolve(event);
};
