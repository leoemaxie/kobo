import { redirect, type Handle } from '@sveltejs/kit';
import { validateSession } from '$lib/server/auth/session';

const PUBLIC_ROUTES = ['/auth/login', '/auth/signup', '/auth/verify-email', '/auth/forgot-password', '/auth/reset-password'];
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
    redirect(302, '/auth/login');
  }

  // If logged in but email not verified, restrict access to verify-email only
  if (
    event.locals.user &&
    !event.locals.user.emailVerifiedAt &&
    path !== '/auth/verify-email' &&
    !isPublic
  ) {
    redirect(302, '/auth/verify-email');
  }

  // Superadmin guard
  if (path.startsWith(SUPERADMIN_PREFIX) && event.locals.user?.role !== 'superadmin') {
    return new Response('Not found', { status: 404 });
  }

  return resolve(event);
};
