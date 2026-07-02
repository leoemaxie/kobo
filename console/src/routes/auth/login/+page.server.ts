import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		// Mock login action
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');

		// In a real app, we would validate credentials and set the session cookie here
		// For now, redirect to the dashboard
		throw redirect(303, '/dashboard');
	}
};
