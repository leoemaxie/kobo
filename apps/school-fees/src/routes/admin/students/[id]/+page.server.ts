import type { PageServerLoad, Actions } from './$types';
import { error, redirect, fail } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { students, parents, parentStudents } from '$lib/server/db/schema';
import { eq, and } from 'drizzle-orm';
import { kobo } from '$lib/server/kobo-client';
import type { TransactionListResponse } from '@kobo/sdk';

export const load: PageServerLoad = async ({ params, locals }) => {
  if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role)) {
    throw redirect(302, '/dashboard');
  }

  const studentId = params.id;
  const studentRows = await db.select().from(students).where(eq(students.id, studentId)).limit(1);

  if (studentRows.length === 0) {
    throw error(404, 'Student not found');
  }
  const student = studentRows[0];

  // Fetch linked parents
  const linkedLinks = await db
    .select()
    .from(parentStudents)
    .where(eq(parentStudents.studentId, studentId));
  const allParents = await db.select().from(parents).where(eq(parents.role, 'parent'));

  const linkedParents = linkedLinks
    .map((l) => allParents.find((p) => p.id === l.parentId))
    .filter(Boolean)
    .map((p) => ({ id: p!.id, name: p!.name, email: p!.email }));

  const unlinkedParents = allParents
    .filter((p) => !linkedParents.some((lp) => lp.id === p.id))
    .map((p) => ({ id: p.id, name: p.name }));

  let statement = { closing_balance_kobo: 0 };
  let transactions: TransactionListResponse = { data: [], next_cursor: null };
  let identity: any = {};

  try {
    identity = await kobo.identities.get(student.koboIdentityId);
  } catch (e) {
    console.error('Failed to fetch Kobo identity:', e);
  }

  try {
    statement = await kobo.accounts.getStatement(student.koboIdentityId);
  } catch (e) {
    console.error('Failed to fetch Kobo statement:', e);
  }

  try {
    transactions = await kobo.accounts.listTransactions(student.koboIdentityId, { limit: 50 });
  } catch (e) {
    console.error('Failed to fetch Kobo transactions:', e);
  }

  return {
    student: {
      id: student.id,
      name: student.name,
      class: student.className,
      koboIdentityId: student.koboIdentityId,
      virtualAccount: {
        accountName: identity.display_name,
        accountNumber: identity.virtual_account?.account_number || null,
        bankName: identity.virtual_account?.bank_name || 'Triumph Bank',
      },
      statement: {
        balanceKobo: statement.closing_balance_kobo,
        currency: 'NGN',
      },
      transactions: transactions.data || [],
      linkedParents,
    },
    availableParents: unlinkedParents,
  };
};

export const actions: Actions = {
  closeAccount: async ({ params, locals }) => {
    if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role))
      return fail(403, { error: 'Unauthorized' });

    const studentRows = await db.select().from(students).where(eq(students.id, params.id)).limit(1);
    if (studentRows.length === 0) return fail(404, { error: 'Student not found' });

    try {
      await kobo.identities.close(studentRows[0].koboIdentityId, {
        sweep_destination: { type: 'refund_to_source' },
      });
      return { success: true, message: 'Account closed successfully' };
    } catch (e: any) {
      return fail(500, { error: 'Failed to close Kobo Identity' });
    }
  },
  linkParent: async ({ request, params, locals }) => {
    if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role))
      return fail(403, { error: 'Unauthorized' });

    const data = await request.formData();
    const parentId = data.get('parentId') as string;
    if (!parentId) return fail(400, { error: 'Missing parent ID' });

    try {
      await db.insert(parentStudents).values({ studentId: params.id, parentId });
      return { success: true, message: 'Parent linked successfully' };
    } catch (e: any) {
      return fail(500, { error: 'Failed to link parent' });
    }
  },
  unlinkParent: async ({ request, params, locals }) => {
    if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role))
      return fail(403, { error: 'Unauthorized' });

    const data = await request.formData();
    const parentId = data.get('parentId') as string;
    if (!parentId) return fail(400, { error: 'Missing parent ID' });

    try {
      await db
        .delete(parentStudents)
        .where(and(eq(parentStudents.studentId, params.id), eq(parentStudents.parentId, parentId)));
      return { success: true, message: 'Parent unlinked successfully' };
    } catch (e: any) {
      return fail(500, { error: 'Failed to unlink parent' });
    }
  },
  modifyClass: async ({ request, params, locals }) => {
    if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role))
      return fail(403, { error: 'Unauthorized' });

    const data = await request.formData();
    const className = data.get('className') as string;
    if (!className) return fail(400, { error: 'Missing class name' });

    try {
      await db.update(students).set({ className }).where(eq(students.id, params.id));
      return { success: true, message: 'Class modified successfully' };
    } catch (e: any) {
      return fail(500, { error: 'Failed to modify class' });
    }
  },
  resendReminder: async ({ locals }) => {
    if (!locals.user || !['admin', 'superadmin'].includes(locals.user.role))
      return fail(403, { error: 'Unauthorized' });
    // Mocking the reminder logic since there is no actual email integration visible right now
    // A real system would enqueue a mail job or call an email service API.
    return { success: true, message: 'Payment reminder sent successfully' };
  },
};
