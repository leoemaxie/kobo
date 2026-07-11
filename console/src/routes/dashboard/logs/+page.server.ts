import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, cookies, locals }) => {
  if (!locals.user) throw redirect(302, '/auth/login');

  const page = url.searchParams.get('page') || '1';
  const method = url.searchParams.get('method') || '';
  const statusCode = url.searchParams.get('status_code') || '';

  const token = cookies.get('session');
  const headers = {
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  const targetUrl = new URL(`${env.CORE_URL}/console/logs`);
  targetUrl.searchParams.set('page', page);
  targetUrl.searchParams.set('limit', '50');
  if (method) targetUrl.searchParams.set('method', method);
  if (statusCode) targetUrl.searchParams.set('status_code', statusCode);

  let logs = [];
  let meta = { total: 0, page: 1, limit: 50, totalPages: 1 };

  try {
    const res = await fetch(targetUrl.toString(), { headers });
    if (res.ok) {
      const payload = await res.json();
      logs = payload.data || [];
      meta = payload.meta || meta;
    } else {
      console.error(`Failed to fetch logs: ${res.status}`);
    }
  } catch (err) {
    console.error('Error fetching paginated logs:', err);
  }

  return {
    paginatedLogs: logs,
    meta,
    filters: { method, statusCode },
  };
};
