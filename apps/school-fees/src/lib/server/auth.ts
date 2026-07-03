import { db } from './db';
import { sessions, parents } from './db/schema';
import { eq } from 'drizzle-orm';

export async function createSession(userId: string) {
    const sessionId = globalThis.crypto.randomUUID().replace(/-/g, '');
    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 7); // 7 days

    await db.insert(sessions).values({
        id: sessionId,
        userId,
        expiresAt
    });

    return sessionId;
}

export async function validateSessionToken(token: string) {
    const result = await db
        .select({ user: parents, session: sessions })
        .from(sessions)
        .innerJoin(parents, eq(sessions.userId, parents.id))
        .where(eq(sessions.id, token))
        .limit(1);

    if (result.length === 0) return { session: null, user: null };

    const { user, session } = result[0];
    if (new Date() >= session.expiresAt) {
        await db.delete(sessions).where(eq(sessions.id, session.id));
        return { session: null, user: null };
    }

    return { session, user };
}
