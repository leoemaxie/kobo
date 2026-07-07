import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { sessions } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';

export const actions: Actions = {
    default: async ({ cookies }) => {
        const sessionId = cookies.get('session');
        if (sessionId) {
            await db.delete(sessions).where(eq(sessions.id, sessionId));
            cookies.delete('session', { path: '/' });
        }
        throw redirect(302, '/login');
    }
};
