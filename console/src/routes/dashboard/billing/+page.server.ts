import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';
import { db } from '$lib/server/db';
import { billingRecords, invoices, apiIntegrators } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect, fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import { withCache } from '$lib/utils/cache';

export const load: PageServerLoad = async ({ locals, setHeaders }) => {
	withCache(setHeaders);

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

	const integrator = await db.query.apiIntegrators.findFirst({
		where: eq(apiIntegrators.id, user.integratorId)
	});
	const plan = integrator?.plan || 'pay_as_you_go';

	const currentRecord = dbBilling[0] || {
		period: new Date().toISOString().slice(0, 7), // e.g. "2026-10"
		accountsProvisioned: 0,
		transactionsProcessed: 0,
		webhookDeliveries: 0,
		amountDueKobo: 0
	};

	const accountsCost = currentRecord.accountsProvisioned * 50;
	const transactionsCost = currentRecord.transactionsProcessed * 2;
	const webhooksCost = currentRecord.webhookDeliveries * 0.5;
	const totalCost = accountsCost + transactionsCost + webhooksCost;

	const getPct = (cost: number) => totalCost === 0 ? 0 : Math.round((cost / totalCost) * 100);

	const nextInvoiceDate = new Date();
	nextInvoiceDate.setMonth(nextInvoiceDate.getMonth() + 1);
	nextInvoiceDate.setDate(1);

	const billingOverview = {
		plan: plan,
		nextInvoiceDate: nextInvoiceDate.toISOString().split('T')[0],
		period: currentRecord.period,
		accrued: `₦${(currentRecord.amountDueKobo / 100).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`,
		usageItems: [
			{ key: 'accounts_provisioned', label: 'Virtual Accounts',  calc: `${currentRecord.accountsProvisioned.toLocaleString()} × ₦50`,    amount: `₦${accountsCost.toLocaleString()}`,  pct: getPct(accountsCost) },
			{ key: 'transaction_fees',     label: 'Transaction Fees',  calc: `${currentRecord.transactionsProcessed.toLocaleString()} × ₦2`,   amount: `₦${transactionsCost.toLocaleString()}`,  pct: getPct(transactionsCost) },
			{ key: 'webhook_calls',        label: 'Webhook Deliveries', calc: `${currentRecord.webhookDeliveries.toLocaleString()} × ₦0.5`, amount: `₦${webhooksCost.toLocaleString()}`,  pct: getPct(webhooksCost) },
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
			const res = await fetch(`${env.CORE_URL}/v1/admin/billing/checkout`, {
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
			if (data.checkout_link) {
				throw redirect(303, data.checkout_link);
			}

			return fail(500, { error: 'Invalid response from payment gateway' });
		} catch (e) {
			if ((e as any).status === 303) throw e;
			return fail(500, { error: 'Failed to connect to payment gateway' });
		}
	}
};
