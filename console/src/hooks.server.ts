import { redirect, type Handle, type HandleServerError } from '@sveltejs/kit';
import { validateSession } from '$lib/server/auth/session';
import { dev } from '$app/environment';

const PUBLIC_ROUTES = ['/auth/login', '/auth/signup', '/auth/verify-email', '/auth/forgot-password', '/auth/reset-password'];
const SUPERADMIN_PREFIX = '/admin';

export const handleError: HandleServerError = ({ error, event }) => {
  const err = error as any;
  console.error(`[handleError] ${event.request.method} ${event.url.pathname}`);
  console.error('Message:', err?.message);
  if (err?.cause) console.error('Root cause:', err.cause);
  return { message: 'Internal server error' };
};

export const handle: Handle = async ({ event, resolve }) => {
  const sessionId = event.cookies.get('session');
  let result = null;
  if (sessionId) {
    try {
      result = await validateSession(sessionId);
    } catch (err: any) {
      console.error('[session] validateSession failed — cause:', err?.cause ?? err);
      // Treat as unauthenticated; clear the broken cookie
      event.cookies.delete('session', { path: '/' });
    }
  }

  event.locals.user = result?.users ?? null;
  event.locals.session = result?.sessions ?? null;

  const path = event.url.pathname;
  const isPublic = PUBLIC_ROUTES.some((p) => path.startsWith(p));

  // If the route doesn't exist (404), let SvelteKit handle it naturally
  if (event.route.id === null) {
    return resolve(event);
  }

  // If not public route and not logged in, redirect to login
  if (!isPublic && !event.locals.user) {
    throw redirect(302, '/auth/login');
  }

  // If logged in but email not verified, restrict access to verify-email only
  if (
    event.locals.user &&
    !event.locals.user.emailVerifiedAt &&
    path !== '/auth/verify-email' &&
    !isPublic
  ) {
    throw redirect(302, '/auth/verify-email');
  }

  if (path.startsWith(SUPERADMIN_PREFIX) && event.locals.user?.role !== 'superadmin') {
    return new Response('Not found', { status: 404 });
  }

  if (path.startsWith('/dashboard') && event.locals.user?.role === 'superadmin') {
    throw redirect(302, '/admin/integrators');
  }

  const startTime = Date.now();
  const response = await resolve(event);
  const duration = Date.now() - startTime;

  if (dev) {
    const isError = response.status >= 400;
    const logPrefix = isError ? '❌' : '✅';
    console.log(`${logPrefix} [${event.request.method}] ${path} - ${response.status} (${duration}ms)`);
  }

  return response;
};
