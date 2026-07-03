import { db } from '$lib/server/db';
import { sessions, users } from '$lib/server/db/schema';
import { eq, and, isNull, gt } from 'drizzle-orm';

const SESSION_DURATION_MS = 1000 * 60 * 60 * 24 * 7; // 7 days

export function generateSessionId(): string {
  const bytes = new Uint8Array(32);
  crypto.getRandomValues(bytes);
  return Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
}

export async function createSession(userId: string) {
  const id = generateSessionId();
  const expiresAt = new Date(Date.now() + SESSION_DURATION_MS);
  await db.insert(sessions).values({ id, userId, expiresAt });
  return { id, expiresAt };
}

export async function validateSession(sessionId: string) {
  const [sessionData] = await db
    .select()
    .from(sessions)
    .where(
      and(eq(sessions.id, sessionId), isNull(sessions.revokedAt), gt(sessions.expiresAt, new Date()))
    )
    .innerJoin(users, eq(users.id, sessions.userId))
    .limit(1);

  return sessionData ?? null; // null means: invalid, expired, or revoked
}

export async function revokeSession(sessionId: string) {
  await db.update(sessions).set({ revokedAt: new Date() }).where(eq(sessions.id, sessionId));
}

// Used by the superadmin "force logout" admin action
export async function revokeAllSessionsForUser(userId: string) {
  await db
    .update(sessions)
    .set({ revokedAt: new Date() })
    .where(and(eq(sessions.userId, userId), isNull(sessions.revokedAt)));
}
