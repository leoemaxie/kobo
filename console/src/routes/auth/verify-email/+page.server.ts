import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { emailVerificationTokens, users } from '$lib/server/db/schema';
import { eq, and, isNull, gt } from 'drizzle-orm';

export const load: PageServerLoad = async ({ url, locals }) => {
	const token = url.searchParams.get('token');

	// If no token is provided in the URL, just return (the +page.svelte tells them to check their email)
	if (!token) {
		if (locals.user?.emailVerifiedAt) {
			throw redirect(302, '/dashboard');
		}
		return {};
	}

	// Validate the token
	const [tokenData] = await db.select()
		.from(emailVerificationTokens)
		.where(
			and(
				eq(emailVerificationTokens.id, token),
				isNull(emailVerificationTokens.usedAt),
				gt(emailVerificationTokens.expiresAt, new Date())
			)
		)
		.limit(1);

	if (!tokenData) {
		return { error: 'Invalid or expired token.' };
	}

	// Mark user as verified and consume the token
	await db.update(users)
		.set({ emailVerifiedAt: new Date(), updatedAt: new Date() })
		.where(eq(users.id, tokenData.userId));

	await db.update(emailVerificationTokens)
		.set({ usedAt: new Date() })
		.where(eq(emailVerificationTokens.id, token));

	// Instantly drop them into the dashboard now that they are verified
	throw redirect(302, '/dashboard');
};
