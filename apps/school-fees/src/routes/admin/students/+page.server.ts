import type { PageServerLoad, Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { students } from '$lib/server/db/schema';
import { kobo } from '$lib/server/kobo-client';
import { eq } from 'drizzle-orm';

export const load: PageServerLoad = async ({ locals }) => {
    if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role)) {
        throw redirect(302, '/dashboard');
    }
    if (locals.user.role === 'admin' && locals.user.status !== 'active') {
        throw redirect(302, '/admin/pending');
    }

    const allStudents = await db.select().from(students).orderBy(students.createdAt);
    
    // Fetch identities from Kobo to get virtual account details
    const koboIdentities = await kobo.identities.list({ limit: 100 });
    const identityMap = new Map(koboIdentities.map(id => [id.id, id]));

    return {
        students: allStudents.map(s => {
            const identity = identityMap.get(s.koboIdentityId);
            return {
                id: s.id,
                name: s.name,
                class: s.className,
                virtualAccountNo: identity?.virtual_account?.account_number || null,
                accountName: identity?.virtual_account?.account_name || null,
                date: s.createdAt.toLocaleDateString(),
                koboIdentityId: s.koboIdentityId
            };
        })
    };
};

export const actions: Actions = {
    register: async ({ request, locals }) => {
        if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role)) {
            return fail(403, { error: 'Unauthorized' });
        }
        if (locals.user.role === 'admin' && locals.user.status !== 'active') {
            return fail(403, { error: 'Account pending approval' });
        }

        const data = await request.formData();
        const name = data.get('name') as string;
        const className = data.get('className') as string;
        const id = data.get('studentId') as string;

        if (!name || !className || !id) {
            return fail(400, { error: 'Missing fields' });
        }

        try {
            const koboResponse = await kobo.identities.create({
                external_reference: id,
                display_name: name,
                metadata: {
                    identity_type: 'individual'
                }
            });

            await db.insert(students).values({
                id,
                name,
                className,
                koboIdentityId: koboResponse.id 
            });

            return { success: true };
        } catch (e: any) {
            return fail(500, { error: 'Failed to provision Kobo Identity: ' + e.message });
        }
    },
    closeAccount: async ({ request, locals }) => {
        if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role)) {
            return fail(403, { error: 'Unauthorized' });
        }
        if (locals.user.role === 'admin' && locals.user.status !== 'active') {
            return fail(403, { error: 'Account pending approval' });
        }

        const data = await request.formData();
        const studentId = data.get('studentId') as string;

        const studentRows = await db.select().from(students).where(eq(students.id, studentId)).limit(1);
        if (studentRows.length === 0) return fail(404, { error: 'Student not found' });

        try {
            await kobo.identities.close(studentRows[0].koboIdentityId, {
                sweep_destination: { type: "refund_to_source" }
            });
            return { success: true };
        } catch (e: any) {
            return fail(500, { error: 'Failed to close Kobo Identity' });
        }
    }
};
