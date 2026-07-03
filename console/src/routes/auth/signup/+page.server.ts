import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { users, apiIntegrators, emailVerificationTokens } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth/session';
import { EmailService } from '$lib/server/email';

export const actions: Actions = {
	default: async ({ request, cookies, url }) => {
		const data = await request.formData();
		const company = data.get('company')?.toString();
		const email = data.get('email')?.toString();
		const password = data.get('password')?.toString();

		if (!company || !email || !password) {
			return fail(400, { error: 'Missing required fields' });
		}

		// Check if user already exists
		const existingUser = await db.query.users.findFirst({
			where: eq(users.email, email)
		});

		if (existingUser) {
			return fail(400, { error: 'Email already exists' });
		}

		// Hash password
		const passwordHash = await argon2.hash(password);

		// Create integrator
		const [integrator] = await db.insert(apiIntegrators)
			.values({ name: company })
			.returning();

		// Create user
		const [user] = await db.insert(users)
			.values({
				integratorId: integrator.id,
				email,
				passwordHash,
				role: 'owner'
			})
			.returning();

		// Generate verification token
		const tokenBytes = new Uint8Array(32);
		crypto.getRandomValues(tokenBytes);
		const token = Array.from(tokenBytes).map(b => b.toString(16).padStart(2, '0')).join('');
		
		await db.insert(emailVerificationTokens).values({
			id: token,
			userId: user.id,
			expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24) // 24 hours
		});

		// Send email asynchronously
		EmailService.sendVerificationEmail(email, token, url.origin).catch(console.error);

		// Create session
		const session = await createSession(user.id);
		cookies.set('session', session.id, {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			expires: session.expiresAt,
			secure: !import.meta.env.DEV
		});

		throw redirect(303, '/auth/verify-email');
	}
};
