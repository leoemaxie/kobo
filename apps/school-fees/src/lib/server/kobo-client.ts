import { createHmac } from 'node:crypto';
import { env } from '$env/dynamic/private';

const KOBO_API_URL = 'https://sandbox.api.kobo.triumphsystems.tech/v1';

export async function koboFetch(
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> {
  const apiKey = env.KOBO_API_KEY;
  const apiSecret = env.KOBO_API_SECRET;

  if (!apiKey || !apiSecret) {
    throw new Error('Missing Kobo API credentials');
  }

  const timestamp = Math.floor(Date.now() / 1000).toString();
  const method = options.method || 'GET';
  const bodyString = options.body ? String(options.body) : '';

  // Signature payload: method + endpoint + timestamp + body
  const payloadToSign = `${method}${endpoint}${timestamp}${bodyString}`;
  const signature = createHmac('sha256', apiSecret).update(payloadToSign).digest('hex');

  const headers = new Headers(options.headers);
  headers.set('Authorization', `Bearer ${apiKey}`);
  headers.set('X-Kobo-Timestamp', timestamp);
  headers.set('X-Kobo-Signature', signature);
  headers.set('Content-Type', 'application/json');

  const url = `${KOBO_API_URL}${endpoint}`;

  return fetch(url, { ...options, headers });
}

export async function createStudentIdentity(reference: string, name: string) {
  const response = await koboFetch('/identities', {
    method: 'POST',
    body: JSON.stringify({
      external_reference: reference,
      display_name: name,
    }),
  });
  if (!response.ok) throw new Error('Failed to create Kobo Identity');
  return response.json();
}

export async function closeStudentIdentity(koboIdentityId: string) {
  const response = await koboFetch(`/identities/${koboIdentityId}/close`, {
    method: 'POST',
    body: JSON.stringify({
      sweep_destination: { type: 'refund_to_source' }
    }),
  });
  if (!response.ok) throw new Error('Failed to close Kobo Identity');
  return response.json();
}

export async function getStudentStatement(koboIdentityId: string) {
  // Assuming identity maps directly or we get account from identity
  // The arch document says: GET /v1/accounts/{accountId}/statement
  // Assuming kobo_identity_id can be used or resolved to accountId.
  const response = await koboFetch(`/accounts/${koboIdentityId}/statement`);
  if (!response.ok) throw new Error('Failed to get statement');
  return response.json();
}

export async function getStudentTransactions(koboIdentityId: string) {
  const response = await koboFetch(`/accounts/${koboIdentityId}/transactions`);
  if (!response.ok) throw new Error('Failed to get transactions');
  return response.json();
}
