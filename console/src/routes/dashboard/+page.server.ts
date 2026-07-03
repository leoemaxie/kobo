import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (!user || !user.integratorId) {
		throw redirect(302, '/auth/login');
	}

	// This data will eventually come from the Kobo Core API (usage stats/metrics)
	// For now, we provide the 1-to-1 mapping expected by the frontend components.
	const metrics = [
		{ key: 'api_requests', label: 'API Requests', value: '1,248', delta: '+12.5%', sub: 'Last 30 days', trend: 'up', bar: 65 },
		{ key: 'virtual_accounts', label: 'Virtual Accounts', value: '342', delta: '+4', sub: '500 limit (sandbox)', trend: 'up', bar: 68 },
		{ key: 'error_rate', label: 'Error Rate', value: '0.8%', delta: '−0.2pp', sub: 'vs. prior period', trend: 'down', bar: 8 },
		{ key: 'p99_latency', label: 'p99 Latency', value: '142ms', delta: '~', sub: 'stable this week', trend: 'neutral', bar: 28 }
	];

	const logs = [
		{ method: 'POST', path: '/v1/accounts', status: 201, ms: 87, id: 'req_9Kz2', time: 'just now' },
		{ method: 'GET', path: '/v1/identities/id_892nf8', status: 200, ms: 43, id: 'req_8mXp', time: '2m ago' },
		{ method: 'POST', path: '/v1/transactions', status: 400, ms: 120, id: 'req_7nQr', time: '15m ago' },
		{ method: 'GET', path: '/v1/accounts?limit=10', status: 200, ms: 31, id: 'req_6wKs', time: '1h ago' },
		{ method: 'DELETE', path: '/v1/hooks/wh_92nd', status: 204, ms: 55, id: 'req_5jVt', time: '3h ago' }
	];

	return {
		metrics,
		logs
	};
};
