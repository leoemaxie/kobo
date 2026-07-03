import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { users, passwordResetTokens } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { EmailService } from '$lib/server/email';

export const actions: Actions = {
	default: async ({ request, url }) => {
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
			const tokenBytes = new Uint8Array(32);
			crypto.getRandomValues(tokenBytes);
			const token = Array.from(tokenBytes).map(b => b.toString(16).padStart(2, '0')).join('');
			
			await db.insert(passwordResetTokens).values({
				id: token,
				userId: user.id,
				expiresAt: new Date(Date.now() + 1000 * 60 * 60) // 1 hour validity
			});

			EmailService.sendPasswordResetEmail(email, token, url.origin).catch(console.error);
		}

		return { success: true };
	}
};
