import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';
import { invoices, apiIntegrators, billingRecords } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { error, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals, params }) => {
  const user = locals.user;
  if (!user || !user.integratorId) {
    throw redirect(302, '/auth/login');
  }

  const invoiceId = params.id;

  const dbInvoice = await db.query.invoices.findFirst({
    where: eq(invoices.id, invoiceId),
  });

  if (!dbInvoice || dbInvoice.integratorId !== user.integratorId) {
    throw error(404, 'Invoice not found');
  }

  const dbIntegrator = await db.query.apiIntegrators.findFirst({
    where: eq(apiIntegrators.id, user.integratorId),
  });

  const dbRecord = await db.query.billingRecords.findFirst({
    where: eq(billingRecords.id, dbInvoice.billingRecordId),
  });

  return {
    invoice: {
      id: dbInvoice.id,
      date: dbInvoice.createdAt.toISOString().split('T')[0],
      period: dbInvoice.period,
      amount: `₦${(dbInvoice.amountKobo / 100).toLocaleString()}`,
      status: dbInvoice.status,
      nombaOrderRef: dbInvoice.nombaOrderRef,
    },
    integrator: dbIntegrator?.name || 'Unknown Workspace',
    details: dbRecord
      ? {
          accounts: dbRecord.accountsProvisioned,
          transactions: dbRecord.transactionsProcessed,
          webhooks: dbRecord.webhookDeliveries,
        }
      : null,
  };
};
