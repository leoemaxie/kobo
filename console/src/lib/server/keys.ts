import { createHash } from "crypto";

export async function generateKeyPair(env: "sandbox" | "production") {
  const prefix = env === "production" ? "kobo_live_" : "kobo_test_";
  const randomBytesId = new Uint8Array(8);
  crypto.getRandomValues(randomBytesId);
  const randomId = Array.from(randomBytesId)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
  const keyId = `${prefix}pk_${randomId}`;

  const randomBytesSecret = new Uint8Array(32);
  crypto.getRandomValues(randomBytesSecret);
  const randomSecretHex = Array.from(randomBytesSecret)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
  const plainSecret = `${prefix}sk_${randomSecretHex}`;

  const secretHash = createHash("sha256").update(plainSecret).digest("hex");

  return { keyId, plainSecret, secretHash };
}
