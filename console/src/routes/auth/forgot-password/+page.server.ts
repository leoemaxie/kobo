import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { users, passwordResetTokens } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { EmailService } from '$lib/server/email';
import { generateToken, hashToken } from '$lib/server/auth/token';

export const actions: Actions = {
	default: async ({ request, url }) => {
		try {
			const data = await request.formData();
			const email = data.get('email')?.toString();

			if (!email) {
				return fail(400, { error: 'Email is required' });
			}

			const user = await db.query.users.findFirst({
				where: eq(users.email, email)
			});

			// Don't leak whether user exists to unauthenticated parties, return success anyway
			if (user) {
				const rawToken = generateToken();
				const tokenHash = hashToken(rawToken);
				
				await db.insert(passwordResetTokens).values({
					id: tokenHash,
					userId: user.id,
					expiresAt: new Date(Date.now() + 1000 * 60 * 60) // 1 hour validity
				});

				EmailService.sendPasswordResetEmail(email, rawToken, url.origin).catch(console.error);
			}

			return { success: true };
		} catch (error) {
			console.error('Forgot password error:', error);
			return fail(500, { error: 'An unexpected error occurred. Please try again later.' });
		}
	}
};
