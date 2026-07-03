import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { apiIntegrators, users } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || user.role !== 'superadmin') {
		throw redirect(302, '/auth/login');
	}

	const allInteg = await db.select({
		id: apiIntegrators.id,
		name: apiIntegrators.name,
		status: apiIntegrators.status,
		prodAccess: apiIntegrators.productionAccessGranted,
		joined: apiIntegrators.createdAt,
		email: users.email
	})
	.from(apiIntegrators)
	.leftJoin(users, eq(apiIntegrators.id, users.integratorId))
	.orderBy(desc(apiIntegrators.createdAt));

	const seen = new Set();
	const mapped = [];
	for (const row of allInteg) {
		if (!seen.has(row.id)) {
			seen.add(row.id);
			mapped.push({
				name: row.name,
				email: row.email || 'No owner email',
				status: row.status,
				prodAccess: row.prodAccess,
				joined: row.joined.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
			});
		}
	}

	return {
		integrators: mapped
	};
};
