import type { PageServerLoad, Actions } from "./$types";
import { env } from "$env/dynamic/private";
import { fail, redirect } from "@sveltejs/kit";
import { db } from "$lib/server/db";
import { apiIntegrators } from "$lib/server/db/schema";
import { eq } from "drizzle-orm";
import { withCache } from "$lib/utils/cache";

export const load: PageServerLoad = async ({
  locals,
  fetch,
  cookies,
  setHeaders,
}) => {
  withCache(setHeaders);
  const user = locals.user;

  if (!user || user.role !== "owner") {
    throw redirect(302, "/dashboard");
  }

  const token = cookies.get("session");
  const headers = {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };

  try {
    // Fetch banks
    const banksRes = await fetch(`${env.CORE_URL}/console/payouts/banks`, {
      headers,
    });
    if (!banksRes.ok) {
      console.error(`Failed to fetch banks: ${banksRes.status} ${banksRes.statusText} - ${await banksRes.text()}`);
    }
    const banks = banksRes.ok ? await banksRes.json() : [];

    // Fetch current bank account
    const accountRes = await fetch(
      `${env.CORE_URL}/console/payouts/bank-account`,
      { headers },
    );
    if (!accountRes.ok && accountRes.status !== 404) {
      console.error(`Failed to fetch bank account: ${accountRes.status} ${accountRes.statusText} - ${await accountRes.text()}`);
    }
    const bankAccount = accountRes.ok ? await accountRes.json() : null;

    // Fetch payout history
    const historyRes = await fetch(`${env.CORE_URL}/console/payouts/`, {
      headers,
    });
    if (!historyRes.ok) {
      console.error(`Failed to fetch payouts: ${historyRes.status} ${historyRes.statusText} - ${await historyRes.text()}`);
    }
    const payouts = historyRes.ok ? await historyRes.json() : [];

    // Get wallet balance
    const integrator = await db.query.apiIntegrators.findFirst({
      where: eq(apiIntegrators.id, user.integratorId!),
    });

    return {
      banks,
      bankAccount,
      payouts,
      walletBalanceKobo: integrator?.walletBalanceKobo || 0,
    };
  } catch (e) {
    console.error("Failed to load payouts data:", e);
    return { banks: [], bankAccount: null, payouts: [], walletBalanceKobo: 0 };
  }
};

export const actions: Actions = {
  saveBankAccount: async ({ request, cookies }) => {
    const data = await request.formData();
    const accountNumber = data.get("accountNumber")?.toString();
    const bankCode = data.get("bankCode")?.toString();

    if (!accountNumber || !bankCode) {
      return fail(400, { error: "Account number and bank code are required" });
    }

    const token = cookies.get("session");

    // First verify account lookup
    const lookupRes = await fetch(
      `${env.CORE_URL}/console/payouts/bank-account/lookup`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ accountNumber, bankCode }),
      },
    );

    if (!lookupRes.ok) {
      return fail(400, { error: "Invalid bank account details" });
    }

    const lookupData = await lookupRes.json();

    // Save the verified bank account
    const saveRes = await fetch(
      `${env.CORE_URL}/console/payouts/bank-account`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          accountNumber,
          bankCode,
          bankName: data.get("bankName")?.toString() || "",
          accountName: lookupData.accountName,
        }),
      },
    );

    if (!saveRes.ok) {
      return fail(500, { error: "Failed to save bank account" });
    }

    return { success: true };
  },

  requestPayout: async ({ request, cookies }) => {
    const data = await request.formData();
    const amount = parseInt(data.get("amount")?.toString() || "0");

    if (isNaN(amount) || amount < 1000) {
      return fail(400, { error: "Invalid amount (minimum ₦1,000)" });
    }

    const token = cookies.get("session");

    const res = await fetch(`${env.CORE_URL}/console/payouts/request`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ requestedAmountKobo: amount * 100 }),
    });

    if (!res.ok) {
      const errorData = await res.text();
      return fail(400, { error: errorData || "Failed to process payout" });
    }

    return { success: true };
  },
};
