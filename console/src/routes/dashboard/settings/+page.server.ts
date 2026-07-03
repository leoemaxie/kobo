import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { apiIntegrators } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const integrator = await db.query.apiIntegrators.findFirst({
		where: eq(apiIntegrators.id, user.integratorId)
	});

	if (!integrator) {
		throw redirect(302, '/auth/login');
	}

	const settings = {
		workspaceName: integrator.name,
		supportEmail: user.email
	};

	return { settings };
};
