import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { users, apiIntegrators } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user) throw redirect(302, '/auth/login');

	// Already has a workspace — skip onboarding
	if (user.integratorId) throw redirect(302, '/dashboard');

	return {};
};

export const actions: Actions = {
	default: async ({ request, locals }) => {
		try {
			const user = locals.user;
			if (!user) return fail(401, { error: 'Unauthorized' });
			if (user.integratorId) throw redirect(302, '/dashboard');

			const data = await request.formData();
			const name = data.get('name')?.toString()?.trim();

			if (!name) {
				return fail(400, { error: 'Workspace name is required' });
			}
			if (name.length < 2) {
				return fail(400, { error: 'Workspace name must be at least 2 characters' });
			}

			// The Core schema requires app-generated UUIDs for api_integrators
			const integratorId = crypto.randomUUID();

			await db.insert(apiIntegrators).values({
				id: integratorId,
				name
			});

			// Link the user to the new workspace
			await db.update(users)
				.set({ integratorId, updatedAt: new Date() })
				.where(eq(users.id, user.id));

			throw redirect(303, '/dashboard');
		} catch (error) {
			if (isRedirect(error)) throw error;
			console.error('Onboarding error:', error);
			return fail(500, { error: 'Failed to create workspace. Please try again.' });
		}
	}
};
