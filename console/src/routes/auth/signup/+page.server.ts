import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/db';
import { users, emailVerificationTokens } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import * as argon2 from 'argon2';
import { createSession } from '$lib/server/auth/session';
import { EmailService } from '$lib/server/email';

export const actions: Actions = {
  default: async ({ request, cookies, url }) => {
    try {
      const data = await request.formData();
      const email = data.get('email')?.toString();
      const password = data.get('password')?.toString();

      if (!email || !password) {
        return fail(400, { error: 'Missing required fields' });
      }

      const existingUser = await db.query.users.findFirst({
        where: eq(users.email, email),
      });

      if (existingUser) {
        return fail(400, {
          error: 'An account with this email already exists',
        });
      }

      const passwordHash = await argon2.hash(password);

      // Create the user only — no integrator yet.
      // The workspace (integrator + api credentials) is created after
      // email verification during the onboarding flow.
      const [user] = await db
        .insert(users)
        .values({
          email,
          passwordHash,
          role: 'owner',
        })
        .returning();

      const tokenBytes = new Uint8Array(32);
      crypto.getRandomValues(tokenBytes);
      const token = Array.from(tokenBytes)
        .map((b) => b.toString(16).padStart(2, '0'))
        .join('');

      await db.insert(emailVerificationTokens).values({
        id: token,
        userId: user.id,
        expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24), // 24 hours
      });

      EmailService.sendVerificationEmail(email, token, url.origin).catch(console.error);

      const session = await createSession(user.id);
      cookies.set('session', session.token, {
        path: '/',
        httpOnly: true,
        sameSite: 'lax',
        expires: session.expiresAt,
        secure: !import.meta.env.DEV,
      });

      throw redirect(303, '/auth/verify-email');
    } catch (error) {
      if (isRedirect(error)) throw error;
      console.error('Signup error:', error);
      return fail(500, {
        error: 'An unexpected error occurred during signup. Please try again.',
      });
    }
  },
};
