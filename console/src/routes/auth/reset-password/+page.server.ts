import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { users, passwordResetTokens } from '$lib/server/db/schema';
import { eq, and, isNull, gt } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { revokeAllSessionsForUser } from '$lib/server/auth/session';

export const load: PageServerLoad = async ({ url }) => {
	const token = url.searchParams.get('token');
	if (!token) throw redirect(302, '/auth/login');

	const [tokenData] = await db.select()
		.from(passwordResetTokens)
		.where(
			and(
				eq(passwordResetTokens.id, token),
				isNull(passwordResetTokens.usedAt),
				gt(passwordResetTokens.expiresAt, new Date())
			)
		)
		.limit(1);

	if (!tokenData) {
		throw redirect(302, '/auth/login?error=invalid_reset_token');
	}

	return { token };
};

export const actions: Actions = {
	default: async ({ request, url }) => {
		try {
			const data = await request.formData();
			const token = url.searchParams.get('token') || data.get('token')?.toString();
			const password = data.get('password')?.toString();

			if (!token || !password) {
				return fail(400, { error: 'Missing required fields' });
			}

			const [tokenData] = await db.select()
				.from(passwordResetTokens)
				.where(
					and(
						eq(passwordResetTokens.id, token),
						isNull(passwordResetTokens.usedAt),
						gt(passwordResetTokens.expiresAt, new Date())
					)
				)
				.limit(1);

			if (!tokenData) {
				return fail(400, { error: 'Invalid or expired token' });
			}

			const passwordHash = await argon2.hash(password);

			await db.update(users)
				.set({ passwordHash, updatedAt: new Date() })
				.where(eq(users.id, tokenData.userId));

			await db.update(passwordResetTokens)
				.set({ usedAt: new Date() })
				.where(eq(passwordResetTokens.id, token));

			// Revoke all existing sessions so old logins are terminated
			await revokeAllSessionsForUser(tokenData.userId);

			throw redirect(303, '/auth/login?reset=success');
		} catch (error) {
			if (isRedirect(error)) {
				throw error;
			}

			console.error('Password reset error:', error);
			return fail(500, { error: 'An unexpected error occurred during password reset.' });
		}
	}
};
