import { fail, redirect } from "@sveltejs/kit";
import type { PageServerLoad, Actions } from "./$types";
import { db } from "$lib/server/db";
import { invitations, apiIntegrators, users } from "$lib/server/db/schema";
import { eq, and, gt, isNull } from "drizzle-orm";

export const load: PageServerLoad = async ({ url, locals }) => {
  const token = url.searchParams.get("token");
  if (!token) throw redirect(302, "/auth/login");

  const [invite] = await db
    .select({
      id: invitations.id,
      email: invitations.email,
      role: invitations.role,
      workspaceName: apiIntegrators.name,
    })
    .from(invitations)
    .innerJoin(apiIntegrators, eq(apiIntegrators.id, invitations.integratorId))
    .where(
      and(
        eq(invitations.id, token),
        gt(invitations.expiresAt, new Date()),
        isNull(invitations.acceptedAt),
      ),
    )
    .limit(1);

  if (!invite) {
    return { error: "This invitation link is invalid or has expired." };
  }

  const isEmailMismatch = locals.user
    ? locals.user.email.toLowerCase() !== invite.email.toLowerCase()
    : false;

  return {
    invite: {
      id: invite.id,
      role: invite.role,
      workspaceName: invite.workspaceName,
    },
    isEmailMismatch,
    user: locals.user,
  };
};

export const actions: Actions = {
  default: async ({ request, locals }) => {
    const user = locals.user;
    if (!user)
      return fail(401, {
        error: "You must be logged in to accept an invitation.",
      });

    const data = await request.formData();
    const token = data.get("token")?.toString();

    if (!token) return fail(400, { error: "Missing token" });

    if (user.integratorId) {
      return fail(400, {
        error:
          "You are already in a workspace. Please leave your current workspace before joining a new one.",
      });
    }

    const [invite] = await db
      .select()
      .from(invitations)
      .where(
        and(
          eq(invitations.id, token),
          gt(invitations.expiresAt, new Date()),
          isNull(invitations.acceptedAt),
        ),
      )
      .limit(1);

    if (!invite)
      return fail(400, { error: "Invalid or expired invitation token." });
    if (invite.email.toLowerCase() !== user.email.toLowerCase()) {
      return fail(400, {
        error: "This invitation was sent to a different email address.",
      });
    }

    await db
      .update(invitations)
      .set({ acceptedAt: new Date() })
      .where(eq(invitations.id, token));
    await db
      .update(users)
      .set({ integratorId: invite.integratorId, role: invite.role as any })
      .where(eq(users.id, user.id));

    throw redirect(303, "/dashboard");
  },
};
