import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		// Mock forgot password action
		const data = await request.formData();
		const email = data.get('email');

		// The svelte component handles the UI state change natively via preventDefault
		// But providing this action prevents 405 errors if JS fails or progressive enhancement is used
		return { success: true };
	}
};
