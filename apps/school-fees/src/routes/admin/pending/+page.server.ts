import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
  if (!locals.user) {
    throw redirect(302, '/login');
  }

  if (locals.user.role === 'admin' && locals.user.status === 'active') {
    throw redirect(302, '/admin/students');
  }

  if (locals.user.role === 'superadmin') {
    throw redirect(302, '/admin/super');
  }

  if (locals.user.role === 'parent') {
    throw redirect(302, '/dashboard');
  }

  return {};
};
