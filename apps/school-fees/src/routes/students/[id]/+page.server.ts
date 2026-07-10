import type { PageServerLoad } from './$types';
import { error, redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { students, parentStudents } from '$lib/server/db/schema';
import { and, eq } from 'drizzle-orm';
import { kobo } from '$lib/server/kobo-client';

export const load: PageServerLoad = async ({ params, locals }) => {
    if (!locals.user) {
        throw redirect(302, '/login');
    }

    const result = await db
        .select({ student: students })
        .from(students)
        .innerJoin(parentStudents, eq(students.id, parentStudents.studentId))
        .where(
            and(
                eq(students.id, params.id),
                eq(parentStudents.parentId, locals.user.id)
            )
        )
        .limit(1);

    if (result.length === 0) {
        throw error(404, 'Student not found or access denied');
    }

    const student = result[0].student;

    let statement;
    let transactions;
    let identity;

    try {
        identity = await kobo.identities.get(student.koboIdentityId);
        statement = await kobo.accounts.getStatement(student.koboIdentityId);
        transactions = await kobo.accounts.listTransactions(student.koboIdentityId);
    } catch (e) {
        console.error("Failed to fetch Kobo data:", e);
        throw error(500, 'Failed to fetch financial records from Kobo');
    }

    return {
        student: {
            id: student.id,
            name: student.name,
            class: student.className,
            virtualAccount: {
                accountName: identity.display_name,
                accountNumber: identity.virtual_account?.account_number || 'Pending',
                bankName: identity.virtual_account?.bank_name || 'Kobo Demo Bank',
            },
            statement: {
                balance: statement.closing_balance_kobo,
                currency: 'NGN',
            },
            transactions: transactions.data || []
        }
    };
};
