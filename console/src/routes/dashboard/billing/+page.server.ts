import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { billingRecords, invoices } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect, fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	const dbBilling = await db.query.billingRecords.findMany({
		where: eq(billingRecords.integratorId, user.integratorId),
		orderBy: [desc(billingRecords.syncedAt)]
	});

	const dbInvoices = await db.query.invoices.findMany({
		where: eq(invoices.integratorId, user.integratorId),
		orderBy: [desc(invoices.createdAt)]
	});

	const mappedInvoices = dbInvoices.map((inv) => ({
		id: inv.id,
		date: inv.createdAt.toISOString().split('T')[0],
		period: inv.period,
		amount: `₦${(inv.amountKobo / 100).toLocaleString()}`,
		status: inv.status
	}));

	const billingOverview = {
		plan: 'pay_as_you_go',
		nextInvoiceDate: '2026-11-01',
		period: 'Oct 1 – Oct 31, 2026',
		accrued: '₦16,355',
		usageItems: [
			{ key: 'accounts_provisioned', label: 'Virtual Accounts',  calc: '152 × ₦50',    amount: '₦7,600',  pct: 45 },
			{ key: 'transaction_fees',     label: 'Transaction Fees',  calc: '3,325 × ₦2',   amount: '₦6,650',  pct: 35 },
			{ key: 'webhook_calls',        label: 'Webhook Deliveries', calc: '8,210 × ₦0.5', amount: '₦4,105',  pct: 20 },
		],
		planDetails: [
			{ key: 'Provisioning Fee', value: '₦50 / account' },
			{ key: 'Transaction Fee', value: '₦2 / transfer' },
			{ key: 'Platform Fee', value: 'None' },
			{ key: 'Support', value: 'Standard email' }
		]
	};

	return { invoices: mappedInvoices, billingOverview };
};

export const actions: Actions = {
	setupPaymentMethod: async ({ locals, url }) => {
		const user = locals.user;
		if (!user || !user.integratorId) return fail(401, { error: 'Unauthorized' });

		const callbackUrl = new URL('/dashboard/billing/callback', url.origin).toString();

		try {
			const res = await fetch(`${url.origin}/v1/admin/billing/checkout`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					integrator_id: user.integratorId,
					type: 'save_card',
					email: user.email,
					callback_url: callbackUrl
				})
			});

			if (!res.ok) {
				const err = await res.text();
				return fail(500, { error: `Payment gateway error: ${err}` });
			}

			const data = await res.json();
			if (data.checkoutLink) {
				throw redirect(303, data.checkoutLink);
			}

			return fail(500, { error: 'Invalid response from payment gateway' });
		} catch (e) {
			if ((e as any).status === 303) throw e;
			return fail(500, { error: 'Failed to connect to payment gateway' });
		}
	}
};
