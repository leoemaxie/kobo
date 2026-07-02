import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		// Mock signup action
		const data = await request.formData();
		const company = data.get('company');
		const email = data.get('email');
		const password = data.get('password');

		// In a real app, we would create the user and redirect to email verification
		// For now, redirect to the dashboard
		throw redirect(303, '/dashboard');
	}
};
