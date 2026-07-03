import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { billingRecords } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const dbBilling = await db.query.billingRecords.findMany({
		where: eq(billingRecords.integratorId, user.integratorId),
		orderBy: [desc(billingRecords.syncedAt)]
	});

	const invoices = dbBilling.map((b) => ({
		id: `inv_${b.period.replace('-', '_')}`,
		date: b.syncedAt.toISOString().split('T')[0],
		period: b.period,
		amount: `₦${(b.amountDueKobo / 100).toLocaleString()}`,
		status: 'paid'
	}));

	return { invoices };
};
