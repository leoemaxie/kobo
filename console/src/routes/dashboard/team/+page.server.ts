import type { PageServerLoad, Actions } from './$types';
import { db } from '$lib/server/db';
import { users, invitations } from '$lib/server/db/schema';
import { eq, desc, and, isNull, gt } from 'drizzle-orm';
import { redirect, fail } from '@sveltejs/kit';
import { withCache } from '$lib/utils/cache';

export const load: PageServerLoad = async ({ locals, setHeaders }) => {
	withCache(setHeaders);

	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const dbUsers = await db.query.users.findMany({
		where: eq(users.integratorId, user.integratorId),
		orderBy: [desc(users.createdAt)]
	});

	const pendingInvites = await db.query.invitations.findMany({
		where: and(
			eq(invitations.integratorId, user.integratorId),
			isNull(invitations.acceptedAt),
			gt(invitations.expiresAt, new Date())
		),
		orderBy: [desc(invitations.createdAt)]
	});

	const members = [
		...dbUsers.map((u) => ({
			id: u.id,
			email: u.email,
			role: u.role,
			status: u.emailVerifiedAt ? 'Active' : 'Pending',
			mfa: false
		})),
		...pendingInvites.map((inv) => ({
			id: inv.id,
			email: inv.email,
			role: inv.role,
			status: 'Invited',
			mfa: false
		}))
	];

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

		const isUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(targetUserId);
		if (isUUID) {
			// Unlink the user from the workspace instead of hard deleting.
			// This preserves their account, audit logs, and API credential attribution.
			await db.update(users).set({ integratorId: null, role: 'owner' }).where(eq(users.id, targetUserId));
		} else {
			await db.delete(invitations).where(eq(invitations.id, targetUserId));
		}

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

		const isUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(targetUserId);
		if (isUUID) {
			await db.update(users).set({ role: newRole }).where(eq(users.id, targetUserId));
		} else {
			await db.update(invitations).set({ role: newRole }).where(eq(invitations.id, targetUserId));
		}

		return { success: true };
	}
};
