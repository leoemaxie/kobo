import { db } from '$lib/server/db';
import { sessions, users } from '$lib/server/db/schema';
import { eq, and, gt } from 'drizzle-orm';
import { generateToken, hashToken } from './token';

const SESSION_DURATION_MS = 1000 * 60 * 60 * 24 * 7; // 7 days

export async function createSession(userId: string) {
  const token = generateToken();
  const sessionId = hashToken(token);
  const expiresAt = new Date(Date.now() + SESSION_DURATION_MS);
  await db.insert(sessions).values({ id: sessionId, userId, expiresAt });
  return { token, expiresAt };
}

export async function validateSession(token: string) {
  const sessionId = hashToken(token);
  const [sessionData] = await db
    .select()
    .from(sessions)
    .where(and(eq(sessions.id, sessionId), gt(sessions.expiresAt, new Date())))
    .innerJoin(users, eq(users.id, sessions.userId))
    .limit(1);

  return sessionData ?? null; // null means: invalid or expired
}

export async function revokeSession(token: string) {
  const sessionId = hashToken(token);
  await db.delete(sessions).where(eq(sessions.id, sessionId));
}

// Used by the superadmin "force logout" admin action
export async function revokeAllSessionsForUser(userId: string) {
  await db.delete(sessions).where(eq(sessions.userId, userId));
}
