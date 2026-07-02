import type { Handle } from '@sveltejs/kit';
import { validateSession } from '$lib/server/auth/session';

const PUBLIC_ROUTES = ['/login', '/signup', '/verify-email'];
const SUPERADMIN_PREFIX = '/admin';

export const handle: Handle = async ({ event, resolve }) => {
  const sessionId = event.cookies.get('session');
  const result = sessionId ? await validateSession(sessionId) : null;

  event.locals.user = result?.users ?? null;
  event.locals.session = result?.sessions ?? null;

  const path = event.url.pathname;
  const isPublic = PUBLIC_ROUTES.some((p) => path.startsWith(p));

  // If not public route and not logged in, redirect to login
  if (!isPublic && !event.locals.user) {
    return Response.redirect(new URL('/login', event.url), 302);
  }

  // If logged in but email not verified, restrict access to verify-email only
  if (
    event.locals.user &&
    !event.locals.user.emailVerifiedAt &&
    path !== '/verify-email' &&
    !isPublic
  ) {
    return Response.redirect(new URL('/verify-email', event.url), 302);
  }

  // Superadmin guard
  if (path.startsWith(SUPERADMIN_PREFIX) && event.locals.user?.role !== 'superadmin') {
    return new Response('Not found', { status: 404 });
  }

  return resolve(event);
};
