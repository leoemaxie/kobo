import { redirect, fail } from "@sveltejs/kit";
import type { PageServerLoad, Actions } from "./$types";
import { db } from "$lib/server/db";
import { emailVerificationTokens, users } from "$lib/server/db/schema";
import { eq, and, isNull, gt } from "drizzle-orm";
import { EmailService } from "$lib/server/email";

export const load: PageServerLoad = async ({ url, locals }) => {
  try {
    const token = url.searchParams.get("token");

    if (!token) {
      if (locals.user?.emailVerifiedAt) throw redirect(302, "/dashboard");
      return { email: locals.user?.email };
    }

    const [tokenData] = await db
      .select()
      .from(emailVerificationTokens)
      .where(
        and(
          eq(emailVerificationTokens.id, token),
          isNull(emailVerificationTokens.usedAt),
          gt(emailVerificationTokens.expiresAt, new Date()),
        ),
      )
      .limit(1);

    if (!tokenData) {
      return { error: "This verification link is invalid or has expired." };
    }

    await db
      .update(users)
      .set({ emailVerifiedAt: new Date(), updatedAt: new Date() })
      .where(eq(users.id, tokenData.userId));

    await db
      .update(emailVerificationTokens)
      .set({ usedAt: new Date() })
      .where(eq(emailVerificationTokens.id, token));

    throw redirect(302, "/dashboard");
  } catch (error: unknown) {
    const skit = await import("@sveltejs/kit");
    if (skit.isRedirect(error)) throw error;
    console.error("Verify email load error:", error);
    return { error: "Something went wrong. Please try again." };
  }
};

export const actions: Actions = {
  resend: async ({ locals, url }) => {
    try {
      const user = locals.user;
      if (!user) return fail(401, { error: "You must be logged in." });
      if (user.emailVerifiedAt)
        return fail(400, { error: "Email is already verified." });

      await db
        .update(emailVerificationTokens)
        .set({ usedAt: new Date() })
        .where(
          and(
            eq(emailVerificationTokens.userId, user.id),
            isNull(emailVerificationTokens.usedAt),
          ),
        );

      const tokenBytes = new Uint8Array(32);
      crypto.getRandomValues(tokenBytes);
      const token = Array.from(tokenBytes)
        .map((b) => b.toString(16).padStart(2, "0"))
        .join("");

      await db.insert(emailVerificationTokens).values({
        id: token,
        userId: user.id,
        expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24),
      });

      EmailService.sendVerificationEmail(user.email, token, url.origin).catch(
        console.error,
      );

      return { success: true };
    } catch (error) {
      console.error("Resend verification error:", error);
      return fail(500, { error: "Failed to resend. Please try again." });
    }
  },
};
