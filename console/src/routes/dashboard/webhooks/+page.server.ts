import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { webhooks } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const dbWebhooks = await db.query.webhooks.findMany({
		where: eq(webhooks.integratorId, user.integratorId),
		orderBy: [desc(webhooks.createdAt)]
	});

	const endpoints = dbWebhooks.map((w) => ({
		url: w.url,
		status: w.status,
		events: w.events as string[],
		secret: w.secret
	}));

	return { endpoints };
};
