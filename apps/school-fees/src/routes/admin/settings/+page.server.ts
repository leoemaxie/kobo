import type { Actions } from './$types';
import { fail } from '@sveltejs/kit';
import { setKoboCredentials } from '$lib/server/kobo-client';

import { env } from '$env/dynamic/private';

export const load = async () => {
  return {
    consoleUrl: env.KOBO_CONSOLE_URL || 'https://console.kobo.triumphsystems.tech',
  };
};

export const actions: Actions = {
  default: async ({ request, locals }) => {
    if (locals.user?.role !== 'superadmin' && locals.user?.role !== 'admin') {
      return fail(403, { error: 'Unauthorized' });
    }

    const data = await request.formData();
    const apiKey = data.get('apiKey') as string;
    const apiSecret = data.get('apiSecret') as string;

    if (!apiKey || !apiSecret) {
      return fail(400, { error: 'Both API Key and Secret are required' });
    }

    setKoboCredentials(apiKey, apiSecret);

    return { success: true };
  },
};
