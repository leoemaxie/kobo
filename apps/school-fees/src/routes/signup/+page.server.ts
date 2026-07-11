import type { Actions } from './$types';
import { fail, redirect, isRedirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parents } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth';
import { dev } from '$app/environment';
import { sendEmail } from '$lib/server/email';

export const actions: Actions = {
  default: async ({ request, cookies }) => {
    try {
      const data = await request.formData();
      const name = data.get('name') as string;
      const email = data.get('email') as string;
      const password = data.get('password') as string;

      if (!name || !email || !password) {
        return fail(400, { error: 'Missing fields' });
      }

      const existing = await db.select().from(parents).where(eq(parents.email, email)).limit(1);
      if (existing.length > 0) {
        return fail(400, { error: 'Email already in use' });
      }

      const passwordHash = await argon2.hash(password);
      const id = globalThis.crypto.randomUUID();

      const roleInput = data.get('role') as string;
      const isAdmin = roleInput === 'admin';
      const role = isAdmin ? 'admin' : 'parent';
      const status = isAdmin ? 'pending' : 'active';

      await db.insert(parents).values({
        id,
        name,
        email,
        passwordHash,
        role,
        status,
        scope: isAdmin ? 'read' : 'full',
      });

      if (isAdmin) {
        // Find superadmin
        const superadmins = await db
          .select()
          .from(parents)
          .where(eq(parents.role, 'superadmin'))
          .limit(1);
        if (superadmins.length > 0) {
          const superadmin = superadmins[0];
          await sendEmail({
            to: superadmin.email,
            subject: 'New Admin Registration Request',
            html: `
                            <div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
                                <h2>New Admin Registration</h2>
                                <p>A new admin has registered and is pending approval:</p>
                                <ul>
                                    <li><strong>Name:</strong> ${name}</li>
                                    <li><strong>Email:</strong> ${email}</li>
                                </ul>
                                <p>Please log in to the admin console to grant or revoke access.</p>
                            </div>
                        `,
          });
        }
      }

      const token = await createSession(id);
      cookies.set('session', token, {
        path: '/',
        httpOnly: true,
        sameSite: 'lax',
        secure: !dev,
        maxAge: 60 * 60 * 24 * 7,
      });

      if (isAdmin) {
        throw redirect(302, '/admin/pending');
      }
      throw redirect(302, '/dashboard');
    } catch (e) {
      if (isRedirect(e)) throw e;
      console.error('Signup action error:', e);
      return fail(500, { error: 'An unexpected error occurred during registration' });
    }
  },
};
