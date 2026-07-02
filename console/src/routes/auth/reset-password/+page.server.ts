import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		// Mock reset password action
		const data = await request.formData();
		const password = data.get('password');

		// In a real app, validate token and update password
		throw redirect(303, '/auth/login');
	}
};
