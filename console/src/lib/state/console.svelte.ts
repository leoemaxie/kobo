import { setContext, getContext } from 'svelte';

const CONSOLE_STATE_KEY = Symbol('CONSOLE_STATE');

export class ConsoleState {
	apiKeys = $state<any[]>([]);
	teamMembers = $state<any[]>([]);
	webhooks = $state<any[]>([]);
	billingInvoices = $state<any[]>([]);
	billingOverview = $state<any>(null); // To cover the plan/period UI hardcoded data
	metrics = $state<any[]>([]);
	logs = $state<any[]>([]); // API Logs from the dashboard
	settings = $state<any>(null);
	user = $state<any>(null);

	// Admin-specific state
	adminIntegrators = $state<any[]>([]);
	adminAuditLogs = $state<any[]>([]);

	// A thin hydrator method that consumes SvelteKit +page.server.ts load data
	hydrate(data: any) {
		if (!data) return;
		
		if (data.keys) this.apiKeys = data.keys;
		if (data.members) this.teamMembers = data.members;
		if (data.endpoints) this.webhooks = data.endpoints;
		if (data.invoices) this.billingInvoices = data.invoices;
		if (data.billingOverview) this.billingOverview = data.billingOverview;
		if (data.settings) this.settings = data.settings;
		if (data.metrics) this.metrics = data.metrics;
		
		// Differentiate between API logs (dashboard) and Audit logs (admin)
		if (data.logs) {
			if (data.logs[0]?.method) this.logs = data.logs; // It's an API log
			else this.adminAuditLogs = data.logs; // It's an Audit log
		}
		
		if (data.user) this.user = data.user;
		if (data.integrators) this.adminIntegrators = data.integrators;
	}
}

/**
 * Initializes the console state and binds it to the Svelte component tree context.
 * Call this exactly once in the root layout or dashboard layout.
 * Prevents SSR state leakage across different user sessions.
 */
export function initConsoleState(): ConsoleState {
	const state = new ConsoleState();
	setContext(CONSOLE_STATE_KEY, state);
	return state;
}

/**
 * Retrieves the reactive console state from context.
 * Safe to call from any deep child component (MetricCards, ApiLogs, etc.)
 */
export function useConsoleState(): ConsoleState {
	const state = getContext<ConsoleState>(CONSOLE_STATE_KEY);
	if (!state) {
		throw new Error('useConsoleState must be used within a component tree initialized with initConsoleState()');
	}
	return state;
}
