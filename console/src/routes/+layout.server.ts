import { db } from '$lib/server/db';
import { apiIntegrators } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';

export const load = async ({ locals }: any) => {
	let integrator = null;
	if (locals.user?.integratorId) {
		const result = await db.select().from(apiIntegrators).where(eq(apiIntegrators.id, locals.user.integratorId)).limit(1);
		integrator = result[0] ?? null;
	}

	return {
		user: locals.user ? { ...locals.user, integrator } : null
	};
};
