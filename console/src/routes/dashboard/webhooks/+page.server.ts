import type { PageServerLoad, Actions, PageServerLoadEvent } from './$types';
import { db } from '$lib/server/db';
import { webhooks } from '$lib/server/db/schema';
import { eq, desc } from 'drizzle-orm';
import { redirect, fail } from '@sveltejs/kit';
import { withCache } from '$lib/utils/cache';

export const load: PageServerLoad = async ({ locals, setHeaders }: PageServerLoadEvent) => {
  withCache(setHeaders);

  const user = locals.user;
  if (!user || !user.integratorId) {
    throw redirect(302, '/auth/login');
  }

  const dbWebhooks = await db.query.webhooks.findMany({
    where: eq(webhooks.integratorId, user.integratorId),
    orderBy: [desc(webhooks.createdAt)],
  });

  const endpoints = dbWebhooks.map((w) => ({
    id: w.id,
    url: w.url,
    status: w.status,
    events: w.events as string[],
    secret: w.secret,
    environment: w.environment,
  }));

  return { endpoints };
};

export const actions: Actions = {
  addEndpoint: async ({ request, locals }) => {
    const user = locals.user;
    if (!user || !user.integratorId) return fail(403, { error: 'Unauthorized' });

    const data = await request.formData();
    const url = data.get('url')?.toString();
    const eventsStr = data.get('events')?.toString();
    const environment = data.get('environment')?.toString() || 'sandbox';

    if (!url || !url.startsWith('https://')) return fail(400, { error: 'Invalid HTTPS URL' });

    const events = eventsStr ? eventsStr.split(',').map((s) => s.trim()) : [];
    if (events.length === 0) return fail(400, { error: 'Must select at least one event' });

    const randomBytes = new Uint8Array(24);
    crypto.getRandomValues(randomBytes);
    const secret = `whsec_${Array.from(randomBytes)
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('')}`;

    await db.insert(webhooks).values({
      integratorId: user.integratorId,
      environment: environment as 'sandbox' | 'production',
      url,
      secret,
      events,
    });

    return { success: true, secret };
  },

  toggleEndpoint: async ({ request, locals }) => {
    const user = locals.user;
    if (!user || !user.integratorId) return fail(403, { error: 'Unauthorized' });

    const data = await request.formData();
    const id = data.get('id')?.toString();
    const currentStatus = data.get('currentStatus')?.toString();
    if (!id || !currentStatus) return fail(400, { error: 'Missing required fields' });

    const newStatus = currentStatus === 'active' ? 'disabled' : 'active';

    await db.update(webhooks).set({ status: newStatus }).where(eq(webhooks.id, id));

    return { success: true };
  },

  deleteEndpoint: async ({ request, locals }) => {
    const user = locals.user;
    if (!user || !user.integratorId) return fail(403, { error: 'Unauthorized' });

    const data = await request.formData();
    const id = data.get('id')?.toString();
    if (!id) return fail(400, { error: 'Missing id' });

    await db.delete(webhooks).where(eq(webhooks.id, id));

    return { success: true };
  },
};
