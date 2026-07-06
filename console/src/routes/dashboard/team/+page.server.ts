import type { PageServerLoad, Actions } from './$types';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect, fail } from '@sveltejs/kit';

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
		id: u.id,
		email: u.email,
		role: u.role,
		status: u.emailVerifiedAt ? 'Active' : 'Pending',
		mfa: false
	}));

	return { members };
};

export const actions: Actions = {
	inviteMember: async ({ request, locals }) => {
		const user = locals.user;
		if (!user || !user.integratorId || user.role !== 'owner') return fail(403, { error: 'Unauthorized' });

		const data = await request.formData();
		const email = data.get('email')?.toString();
		const role = data.get('role')?.toString() as 'member' | 'owner';

		if (!email || !role) return fail(400, { error: 'Missing fields' });

		const existing = await db.query.users.findFirst({
			where: eq(users.email, email)
		});
		if (existing && existing.integratorId === user.integratorId) {
			return fail(400, { error: 'User is already a member' });
		}

		const { createInvitation } = await import('$lib/server/invitations');
		await createInvitation(user.integratorId, user.id, email, role);

		return { success: true };
	},

	removeMember: async ({ request, locals }) => {
		const user = locals.user;
		if (!user || !user.integratorId || user.role !== 'owner') return fail(403, { error: 'Unauthorized' });

		const data = await request.formData();
		const targetUserId = data.get('id')?.toString();
		if (!targetUserId) return fail(400, { error: 'Missing user id' });

		if (targetUserId === user.id) return fail(400, { error: 'Cannot remove yourself' });

		await db.delete(users).where(eq(users.id, targetUserId));

		return { success: true };
	},

	changeRole: async ({ request, locals }) => {
		const user = locals.user;
		if (!user || !user.integratorId || user.role !== 'owner') return fail(403, { error: 'Unauthorized' });

		const data = await request.formData();
		const targetUserId = data.get('id')?.toString();
		const newRole = data.get('role')?.toString() as 'member' | 'owner';
		if (!targetUserId || !newRole) return fail(400, { error: 'Missing fields' });

		if (targetUserId === user.id) return fail(400, { error: 'Cannot change your own role here' });

		await db.update(users).set({ role: newRole }).where(eq(users.id, targetUserId));

		return { success: true };
	}
};
