import { redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const user = locals.user;
  if (!user) {
    throw redirect(302, "/auth/login");
  }

  const orderRef = url.searchParams.get("orderRef");
  if (orderRef) {
    try {
      const res = await fetch(`${env.CORE_URL}/v1/admin/billing/verify`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          integrator_id: user.integratorId,
          order_ref: orderRef,
        }),
      });

      if (!res.ok) {
        throw redirect(
          302,
          "/dashboard/billing?payment_error=verification_failed",
        );
      }

      throw redirect(302, "/dashboard/billing?payment_success=true");
    } catch (e) {
      if ((e as any).status === 302) throw e;
      throw redirect(
        302,
        "/dashboard/billing?payment_error=verification_error",
      );
    }
  }

  throw redirect(302, "/dashboard/billing");
};
