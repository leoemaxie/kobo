import type { Actions } from './$types';
import { fail, redirect, isRedirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parents } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth';
import { dev } from '$app/environment';

export const actions: Actions = {
  default: async ({ request, cookies }) => {
    try {
      const data = await request.formData();
      const email = data.get('email') as string;
      const password = data.get('password') as string;

      if (!email || !password) {
        return fail(400, { error: 'Missing credentials' });
      }

      const userResults = await db.select().from(parents).where(eq(parents.email, email)).limit(1);
      if (userResults.length === 0) {
        return fail(400, { error: 'Invalid credentials' });
      }

      const user = userResults[0];
      const validPassword = await argon2.verify(user.passwordHash, password);

      if (!validPassword) {
        return fail(400, { error: 'Invalid credentials' });
      }

      const token = await createSession(user.id);
      cookies.set('session', token, {
        path: '/',
        httpOnly: true,
        sameSite: 'lax',
        secure: !dev,
        maxAge: 60 * 60 * 24 * 7,
      });

      if (user.role === 'superadmin') {
        throw redirect(302, '/admin/super');
      } else if (user.role === 'admin') {
        if (user.status === 'pending') {
          throw redirect(302, '/admin/pending');
        }
        throw redirect(302, '/admin/students');
      }
      throw redirect(302, '/dashboard');
    } catch (e) {
      if (isRedirect(e)) throw e;
      console.error('Login action error:', e);
      return fail(500, { error: 'An unexpected error occurred during login' });
    }
  },
};
