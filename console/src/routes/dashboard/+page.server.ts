import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { withCache } from '$lib/utils/cache';
import { db } from '$lib/server/db';
import { apiCredentials, webhooks, usageEvents, paymentMethods } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';

import { env } from '$env/dynamic/private';

export const load: PageServerLoad = async ({ locals, setHeaders, cookies, fetch }) => {
  withCache(setHeaders);

  const user = locals.user;
  if (!user) throw redirect(302, '/auth/login');

  // If the user hasn't created a workspace yet, send them to onboarding
  if (!user.integratorId) throw redirect(302, '/dashboard/onboarding');

  const hasKeys = !!(await db.query.apiCredentials.findFirst({
    where: eq(apiCredentials.integratorId, user.integratorId),
  }));

  const hasWebhooks = !!(await db.query.webhooks.findFirst({
    where: eq(webhooks.integratorId, user.integratorId),
  }));

  const hasUsage = !!(await db.query.usageEvents.findFirst({
    where: eq(usageEvents.integratorId, user.integratorId),
  }));

  const hasBilling = !!(await db.query.paymentMethods.findFirst({
    where: eq(paymentMethods.integratorId, user.integratorId),
  }));

  const token = cookies.get('session');
  const headers = {
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  let metrics = [];
  let logs = [];

  try {
    const url = `${env.CORE_URL}/console/analytics`;
    console.log(`Fetching analytics from: ${url}`);
    const res = await fetch(url, { headers });
    console.log(`Analytics response status: ${res.status}`);
    if (res.ok) {
      const data = await res.json();
      console.log(`Analytics data:`, data);
      metrics = data.metrics || [];
      logs = data.logs || [];
    } else {
      console.error(`Failed to fetch analytics: ${res.status} ${res.statusText}`);
      console.error(`Response body: ${await res.text()}`);
    }
  } catch (err) {
    console.error('Error fetching analytics:', err);
  }

  return {
    metrics,
    logs,
    setupStatus: {
      hasKeys,
      hasWebhooks,
      hasUsage,
      hasBilling,
      isProduction: user.integrator?.productionAccessGranted || false,
    },
  };
};
