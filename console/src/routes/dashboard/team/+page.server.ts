import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const dbUsers = await db.query.users.findMany({
		where: eq(users.integratorId, user.integratorId),
		orderBy: [desc(users.createdAt)]
	});

	const members = dbUsers.map((u) => ({
		email: u.email,
		role: u.role,
		status: u.emailVerifiedAt ? 'Active' : 'Pending',
		mfa: false
	}));

	return { members };
};
