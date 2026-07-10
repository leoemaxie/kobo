import { db } from "$lib/server/db";
import { invitations, apiIntegrators } from "$lib/server/db/schema";
import { eq } from "drizzle-orm";
import type { InferInsertModel } from "drizzle-orm";
import { env } from "$env/dynamic/private";
import { EmailService } from "$lib/server/email";

// In a real app, you would send an email here using an email provider like Resend or SendGrid.
export async function createInvitation(
  integratorId: string,
  invitedBy: string,
  email: string,
  role: "member" | "owner" | "superadmin",
) {
  const randomBytesId = new Uint8Array(16);
  crypto.getRandomValues(randomBytesId);
  const token = Array.from(randomBytesId)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");

  // Expires in 7 days
  const expiresAt = new Date(Date.now() + 1000 * 60 * 60 * 24 * 7);

  const invitation: InferInsertModel<typeof invitations> = {
    id: token,
    integratorId,
    invitedBy,
    email,
    role,
    expiresAt,
  };

  await db.insert(invitations).values(invitation);

  const domain = env.KOBO_DOMAIN || "kobo.dev";
  const baseUrl = env.PUBLIC_URL || "https://console.kobo.triumphsystems.tech";

  const workspace = await db.query.apiIntegrators.findFirst({
    where: eq(apiIntegrators.id, integratorId),
  });
  const workspaceName = workspace?.name || "a Kobo workspace";

  await EmailService.sendInvitationEmail(
    email,
    role,
    workspaceName,
    token,
    baseUrl,
  );

  return token;
}
