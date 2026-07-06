import type { PageServerLoad, Actions } from './$types';
import { db } from '$lib/server/db';
import { apiCredentials } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect, fail } from '@sveltejs/kit';

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

export const actions: Actions = {
	createKey: async ({ request, locals }) => {
		const user = locals.user;
		if (!user || !user.integratorId) return fail(403, { error: 'Unauthorized' });

		const data = await request.formData();
		const label = data.get('label')?.toString() || 'New Key';
		const env = data.get('environment')?.toString() as 'sandbox' | 'production' || 'sandbox';
		const scopesStr = data.get('scopes')?.toString();
		const ipsStr = data.get('ips')?.toString();
		
		const scopes = scopesStr ? scopesStr.split(',').map(s => s.trim()) : [];
		const allowedIps = ipsStr ? ipsStr.split('\n').map(s => s.trim()).filter(Boolean) : [];

		const { generateKeyPair } = await import('$lib/server/keys');
		const { keyId, plainSecret, secretHash } = await generateKeyPair(env);

		await db.insert(apiCredentials).values({
			integratorId: user.integratorId,
			environment: env,
			keyId,
			secretHash,
			label,
			createdBy: user.id,
			scopes,
			allowedIps
		});

		return { success: true, plainSecret };
	},

	revokeKey: async ({ request, locals }) => {
		const user = locals.user;
		if (!user || !user.integratorId) return fail(403, { error: 'Unauthorized' });

		const data = await request.formData();
		const keyId = data.get('keyId')?.toString();
		if (!keyId) return fail(400, { error: 'Missing keyId' });

		await db.update(apiCredentials)
			.set({ revokedAt: new Date(), revokedBy: user.id })
			.where(eq(apiCredentials.keyId, keyId));

		return { success: true };
	},

	rollKey: async ({ request, locals }) => {
		const user = locals.user;
		if (!user || !user.integratorId) return fail(403, { error: 'Unauthorized' });

		const data = await request.formData();
		const keyId = data.get('keyId')?.toString();
		if (!keyId) return fail(400, { error: 'Missing keyId' });

		const existing = await db.query.apiCredentials.findFirst({
			where: eq(apiCredentials.keyId, keyId)
		});
		if (!existing) return fail(404, { error: 'Key not found' });

		const { generateKeyPair } = await import('$lib/server/keys');
		const { keyId: newKeyId, plainSecret, secretHash } = await generateKeyPair(existing.environment);

		await db.transaction(async (tx) => {
			await tx.update(apiCredentials)
				.set({ revokedAt: new Date(), revokedBy: user.id, revokedReason: 'rolled' })
				.where(eq(apiCredentials.keyId, keyId));

			await tx.insert(apiCredentials).values({
				integratorId: user.integratorId!,
				environment: existing.environment,
				keyId: newKeyId,
				secretHash,
				label: existing.label,
				createdBy: user.id,
				scopes: existing.scopes,
				allowedIps: existing.allowedIps
			});
		});

		return { success: true, plainSecret };
	}
};
