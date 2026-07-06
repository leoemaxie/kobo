import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals }) => {
	const user = locals.user;
	if (!user) {
		throw redirect(302, '/auth/login');
	}

	// Nomba redirects back to this URL after checkout with query parameters.
	// For instance: ?orderRef=ref_12345
	// We just redirect back to the billing dashboard with a success message.
	
	const orderRef = url.searchParams.get('orderRef');
	if (orderRef) {
		// You might want to verify the orderRef status from Nomba API via Kobo Core
		// but for now we just assume success.
		throw redirect(302, '/dashboard/billing?payment_success=true');
	}

	throw redirect(302, '/dashboard/billing');
};
