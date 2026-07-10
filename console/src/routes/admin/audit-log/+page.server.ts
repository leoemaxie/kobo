import type { PageServerLoad } from "./$types";
import { db } from "$lib/server/db";
import { adminAuditLog, users, apiIntegrators } from "$lib/server/db/schema";
import { eq, desc } from "drizzle-orm";
import { redirect } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ locals }) => {
  const user = locals.user;
  if (!user || user.role !== "superadmin") {
    throw redirect(302, "/auth/login");
  }

  const logs = await db
    .select({
      time: adminAuditLog.createdAt,
      actor: users.email,
      action: adminAuditLog.action,
      target: apiIntegrators.name,
      detail: adminAuditLog.detail,
    })
    .from(adminAuditLog)
    .leftJoin(users, eq(adminAuditLog.actorUserId, users.id))
    .leftJoin(
      apiIntegrators,
      eq(adminAuditLog.targetIntegratorId, apiIntegrators.id),
    )
    .orderBy(desc(adminAuditLog.createdAt))
    .limit(50);

  const mappedLogs = logs.map((l) => ({
    time: l.time.toISOString().replace("T", " ").substring(0, 19) + " UTC",
    actor: l.actor || "Unknown",
    action: l.action,
    target: l.target || "N/A",
    detail: JSON.stringify(l.detail),
  }));

  return {
    logs: mappedLogs,
  };
};
