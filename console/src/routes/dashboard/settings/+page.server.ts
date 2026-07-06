import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { apiIntegrators, users, apiCredentials, billingRecords, webhooks, sessions, emailVerificationTokens, passwordResetTokens, adminAuditLog } from '$lib/server/db/schema';
import { eq, inArray, or } from 'drizzle-orm';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user) throw redirect(302, '/auth/login');

	const integrator = await db.query.apiIntegrators.findFirst({
		where: eq(apiIntegrators.id, user.integratorId!)
	});

	return {
		integrator
	};
};

export const actions: Actions = {
	updateWorkspace: async ({ request, locals }) => {
		try {
			const user = locals.user;
			if (!user || user.role !== 'owner') {
				return fail(403, { error: 'Only owners can update workspace settings' });
			}

			const data = await request.formData();
			const name = data.get('name')?.toString();

			if (!name) {
				return fail(400, { error: 'Missing required fields' });
			}

			await db.update(apiIntegrators)
				.set({ name, updatedAt: new Date() })
				.where(eq(apiIntegrators.id, user.integratorId!));

			return { success: true };
		} catch (error) {
			console.error('Update workspace error:', error);
			return fail(500, { error: 'An unexpected error occurred.' });
		}
	},

	deleteWorkspace: async ({ locals, cookies }) => {
		try {
			const user = locals.user;
			if (!user || user.role !== 'owner') {
				return fail(403, { error: 'Unauthorized' });
			}

			const intId = user.integratorId!;

			// Get all users in the workspace to delete their sessions and tokens
			const workspaceUsers = await db.select({ id: users.id }).from(users).where(eq(users.integratorId, intId));
			const userIds = workspaceUsers.map(u => u.id);

			// Manual cascade delete to avoid Postgres foreign key violations
			await db.delete(webhooks).where(eq(webhooks.integratorId, intId));
			await db.delete(billingRecords).where(eq(billingRecords.integratorId, intId));
			await db.delete(apiCredentials).where(eq(apiCredentials.integratorId, intId));
			
			await db.delete(adminAuditLog).where(or(
				eq(adminAuditLog.targetIntegratorId, intId),
				userIds.length > 0 ? inArray(adminAuditLog.actorUserId, userIds) : undefined,
				userIds.length > 0 ? inArray(adminAuditLog.targetUserId, userIds) : undefined
			));

			if (userIds.length > 0) {
				await db.delete(sessions).where(inArray(sessions.userId, userIds));
				await db.delete(emailVerificationTokens).where(inArray(emailVerificationTokens.userId, userIds));
				await db.delete(passwordResetTokens).where(inArray(passwordResetTokens.userId, userIds));
			}
			
			// Unlink users then delete
			await db.delete(users).where(eq(users.integratorId, intId));
			
			// Finally delete the workspace itself
			await db.delete(apiIntegrators).where(eq(apiIntegrators.id, intId));

			// Kill the session cookie
			cookies.delete('session', { path: '/' });
			throw redirect(303, '/auth/signup');
		} catch (error) {
			import('@sveltejs/kit').then(({ isRedirect }) => {
				if (isRedirect(error)) throw error;
			});
			if ((error as any)?.status === 303) throw error; // Handle SvelteKit redirect error object directly as a fallback if isRedirect isn't available synchronously
			
			console.error('Delete workspace error:', error);
			return fail(500, { error: 'Failed to delete workspace. Please check constraints.' });
		}
	}
};
