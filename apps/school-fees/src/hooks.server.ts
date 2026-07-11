import { validateSessionToken } from '$lib/server/auth';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    const sessionToken = event.cookies.get('session');
    
    if (!sessionToken) {
        event.locals.user = null;
        event.locals.session = null;
    } else {
        const { session, user } = await validateSessionToken(sessionToken);
        if (session) {
            event.locals.user = user;
            event.locals.session = session;
        } else {
            event.locals.user = null;
            event.locals.session = null;
            event.cookies.delete('session', { path: '/' });
        }
    }

    const isPublic = ['/login', '/signup'].includes(event.url.pathname);

    // If root, redirect based on role or to login
    if (event.url.pathname === '/') {
        if (event.locals.user) {
            throw redirect(302, (event.locals.user.role === 'admin' || event.locals.user.role === 'superadmin') ? '/admin/students' : '/dashboard');
        } else {
            throw redirect(302, '/login');
        }
    }

    if (!isPublic && !event.locals.user) {
        throw redirect(302, '/login');
    }

    if (isPublic && event.locals.user) {
        throw redirect(302, (event.locals.user.role === 'admin' || event.locals.user.role === 'superadmin') ? '/admin/students' : '/dashboard');
    }

    return resolve(event);
};
