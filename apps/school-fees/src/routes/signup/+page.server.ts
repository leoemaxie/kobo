import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parents } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth';
import { dev } from '$app/environment';

export const actions: Actions = {
    default: async ({ request, cookies }) => {
        const data = await request.formData();
        const name = data.get('name') as string;
        const email = data.get('email') as string;
        const password = data.get('password') as string;

        if (!name || !email || !password) {
            return fail(400, { error: 'Missing fields' });
        }

        const existing = await db.select().from(parents).where(eq(parents.email, email)).limit(1);
        if (existing.length > 0) {
            return fail(400, { error: 'Email already in use' });
        }

        const passwordHash = await argon2.hash(password);
        const id = globalThis.crypto.randomUUID();

        await db.insert(parents).values({
            id,
            name,
            email,
            passwordHash,
            isAdmin: email.endsWith('@triumph.edu')
        });

        const token = await createSession(id);
        cookies.set('session', token, {
            path: '/',
            httpOnly: true,
            sameSite: 'lax',
            secure: !dev,
            maxAge: 60 * 60 * 24 * 7
        });

        const isAdmin = email.endsWith('@triumph.edu');
        throw redirect(302, isAdmin ? '/admin/students' : '/dashboard');
    }
};
