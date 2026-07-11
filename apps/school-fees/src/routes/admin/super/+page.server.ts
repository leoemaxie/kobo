import type { PageServerLoad, Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parents } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';

export const load: PageServerLoad = async ({ locals }) => {
    if (!locals.user || locals.user.role !== 'superadmin') {
        throw redirect(302, '/dashboard');
    }

    const allAdmins = await db.select().from(parents).where(eq(parents.role, 'admin')).orderBy(parents.createdAt);
    
    return {
        admins: allAdmins.map(a => ({
            id: a.id,
            name: a.name,
            email: a.email,
            status: a.status,
            scope: a.scope || '',
            date: a.createdAt.toLocaleDateString()
        }))
    };
};

export const actions: Actions = {
    updateStatus: async ({ request, locals }) => {
        if (!locals.user || locals.user.role !== 'superadmin') {
            return fail(403, { error: 'Unauthorized' });
        }

        const data = await request.formData();
        const adminId = data.get('adminId') as string;
        const status = data.get('status') as "pending" | "active" | "revoked";

        if (!['pending', 'active', 'revoked'].includes(status)) {
            return fail(400, { error: 'Invalid status' });
        }

        await db.update(parents).set({ status }).where(eq(parents.id, adminId));
        return { success: true };
    },
    updateScope: async ({ request, locals }) => {
        if (!locals.user || locals.user.role !== 'superadmin') {
            return fail(403, { error: 'Unauthorized' });
        }

        const data = await request.formData();
        const adminId = data.get('adminId') as string;
        const scope = data.get('scope') as string;

        await db.update(parents).set({ scope }).where(eq(parents.id, adminId));
        return { success: true };
    }
};
