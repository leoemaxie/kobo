import { fail, redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parents } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth';

export const actions = {
    default: async ({ request, cookies }) => {
        const data = await request.formData();
        const email = data.get('email') as string;
        const password = data.get('password') as string;

        if (!email || !password) {
            return fail(400, { error: 'Missing credentials' });
        }

        const userResults = await db.select().from(parents).where(eq(parents.email, email)).limit(1);
        if (userResults.length === 0) {
            return fail(400, { error: 'Invalid credentials' });
        }

        const user = userResults[0];
        const validPassword = await argon2.verify(user.passwordHash, password);
        
        if (!validPassword) {
            return fail(400, { error: 'Invalid credentials' });
        }

        const token = await createSession(user.id);
        cookies.set('session', token, {
            path: '/',
            httpOnly: true,
            sameSite: 'lax',
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 7
        });

        throw redirect(302, '/dashboard');
    }
};
