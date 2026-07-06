import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth/session';

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		try {
			const data = await request.formData();
			const email = data.get('email')?.toString();
			const password = data.get('password')?.toString();

			if (!email || !password) {
				return fail(400, { error: 'Missing required fields' });
			}

			const user = await db.query.users.findFirst({
				where: eq(users.email, email)
			})

			if (!user) {
				return fail(400, { error: 'Invalid email or password' });
			}

			const validPassword = await argon2.verify(user.passwordHash, password);
			if (!validPassword) {
				return fail(400, { error: 'Invalid email or password' });
			}

			const session = await createSession(user.id);
			cookies.set('session', session.id, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				expires: session.expiresAt,
				secure: !import.meta.env.DEV
			});

			throw redirect(303, '/dashboard');
		} catch (error) {
			if (isRedirect(error)) throw error;
			console.error('Login error:', error);
			return fail(500, { error: 'An unexpected error occurred. Please try again later.' });
		}
	}
};
