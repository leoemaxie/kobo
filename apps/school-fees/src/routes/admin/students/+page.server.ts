import { fail, redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { students } from '$lib/server/db/schema';
import crypto from 'crypto';
import { koboFetch } from '$lib/server/kobo-client';
import { eq } from 'drizzle-orm';

export const load = async ({ locals }) => {
    if (!locals.user || !locals.user.isAdmin) {
        throw redirect(302, '/dashboard');
    }

    const allStudents = await db.select().from(students).orderBy(students.createdAt);
    
    return {
        students: allStudents.map(s => ({
            id: s.id,
            name: s.name,
            class: s.className,
            date: s.createdAt.toLocaleDateString(),
            koboIdentityId: s.koboIdentityId
        }))
    };
};

export const actions = {
    register: async ({ request, locals }) => {
        if (!locals.user || !locals.user.isAdmin) {
            return fail(403, { error: 'Unauthorized' });
        }

        const data = await request.formData();
        const name = data.get('name') as string;
        const className = data.get('className') as string;

        if (!name || !className) {
            return fail(400, { error: 'Missing fields' });
        }

        const id = crypto.randomUUID();

        try {
            const koboResponse = await koboFetch('/identities', {
                method: 'POST',
                body: JSON.stringify({
                    external_reference: id,
                    display_name: name,
                    identity_type: 'individual'
                })
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
        if (!locals.user || !locals.user.isAdmin) {
            return fail(403, { error: 'Unauthorized' });
        }

        const data = await request.formData();
        const studentId = data.get('studentId') as string;

        const studentRows = await db.select().from(students).where(eq(students.id, studentId)).limit(1);
        if (studentRows.length === 0) return fail(404, { error: 'Student not found' });

        try {
            await koboFetch(`/identities/${studentRows[0].koboIdentityId}/close`, {
                method: 'POST',
                body: JSON.stringify({
                    sweep_destination: { type: "refund_to_source" }
                })
            });
            return { success: true };
        } catch (e: any) {
            return fail(500, { error: 'Failed to close Kobo Identity' });
        }
    }
};
