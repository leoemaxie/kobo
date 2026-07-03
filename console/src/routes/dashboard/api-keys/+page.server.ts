import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { apiCredentials } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const dbKeys = await db.query.apiCredentials.findMany({
		where: eq(apiCredentials.integratorId, user.integratorId),
		orderBy: [desc(apiCredentials.createdAt)]
	});

	const keys = dbKeys.map((k) => ({
		name: k.label || 'Default Key',
		id: k.keyId,
		secret: k.secretHash,
		created: k.createdAt.toISOString().split('T')[0],
		lastUsed: 'N/A',
		environment: k.environment,
		status: k.revokedAt ? 'revoked' : 'active'
	}));

	return { keys };
};
