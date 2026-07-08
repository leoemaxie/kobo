import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
  submitKyc: async ({ request }) => {
    const formData = await request.formData();
    const businessName = formData.get('businessName');
    const regNumber = formData.get('regNumber');
    const address = formData.get('address');
    const directorName = formData.get('directorName');
    const bvn = formData.get('bvn');

    if (!businessName || !regNumber || !address || !directorName || !bvn) {
      return fail(400, { error: 'All text fields are required.' });
    }

    // Simulate processing delay for the mock KYC submission
    await new Promise((resolve) => setTimeout(resolve, 1500));
    
    // Normally, documents would be uploaded to a bucket and details stored to a DB here.

    return { success: true };
  }
};
