import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { users, emailVerificationTokens, invitations } from '$lib/server/db/schema';
import { eq, and, gt, isNull } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth/session';
import { EmailService } from '$lib/server/email';

export const actions: Actions = {
	default: async ({ request, cookies, url }) => {
		try {
			const data = await request.formData();
			const email = data.get('email')?.toString();
			const password = data.get('password')?.toString();
			const tokenParam = data.get('token')?.toString();

			if (!email || !password) {
				return fail(400, { error: 'Missing required fields' });
			}

			const existingUser = await db.query.users.findFirst({
				where: eq(users.email, email)
			});

			if (existingUser) {
				return fail(400, { error: 'An account with this email already exists' });
			}

			const passwordHash = await argon2.hash(password);

			let integratorId: string | null = null;
			let role: 'owner' | 'member' | 'superadmin' = 'owner';

			if (tokenParam) {
				const [invite] = await db.select().from(invitations).where(
					and(
						eq(invitations.id, tokenParam),
						gt(invitations.expiresAt, new Date()),
						isNull(invitations.acceptedAt)
					)
				).limit(1);

				if (!invite) {
					return fail(400, { error: 'Invalid or expired invitation token' });
				}

				if (invite.email.toLowerCase() !== email.toLowerCase()) {
					return fail(400, { error: 'This invitation is for a different email address' });
				}

				integratorId = invite.integratorId;
				role = invite.role as 'owner' | 'member' | 'superadmin';

				await db.update(invitations).set({ acceptedAt: new Date() }).where(eq(invitations.id, tokenParam));
			}

			// Create the user
			// If integratorId is null, the workspace (integrator + api credentials) is created after
			// email verification during the onboarding flow. If set, they skip onboarding.
			const [user] = await db.insert(users)
				.values({
					email,
					passwordHash,
					role,
					integratorId
				})
				.returning();

			const tokenBytes = new Uint8Array(32);
			crypto.getRandomValues(tokenBytes);
			const token = Array.from(tokenBytes).map(b => b.toString(16).padStart(2, '0')).join('');

			await db.insert(emailVerificationTokens).values({
				id: token,
				userId: user.id,
				expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24) // 24 hours
			});

			EmailService.sendVerificationEmail(email, token, url.origin).catch(console.error);

			const session = await createSession(user.id);
			cookies.set('session', session.id, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				expires: session.expiresAt,
				secure: !import.meta.env.DEV
			});

			throw redirect(303, '/auth/verify-email');
		} catch (error) {
			if (isRedirect(error)) throw error;
			console.error('Signup error:', error);
			return fail(500, { error: 'An unexpected error occurred during signup. Please try again.' });
		}
	}
};

