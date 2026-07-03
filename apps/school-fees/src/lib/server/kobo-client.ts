import crypto from 'crypto';
import { env } from '$env/dynamic/private';

export async function koboFetch(endpoint: string, options: RequestInit = {}) {
    const KOBO_API_URL = env.KOBO_API_URL || 'http://localhost:8080/v1';
    const apiKey = env.KOBO_API_KEY as string;
    const apiSecret = env.KOBO_API_SECRET as string;

    const timestamp = Math.floor(Date.now() / 1000).toString();
    const bodyString = options.body ? (typeof options.body === 'string' ? options.body : JSON.stringify(options.body)) : '';
    
    // Create signature payload: timestamp + "." + body
    const payload = `${timestamp}.${bodyString}`;
    const hmac = crypto.createHmac('sha256', apiSecret);
    hmac.update(payload);
    const signature = hmac.digest('hex');

    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${apiKey}`,
        'x-kobo-timestamp': timestamp,
        'x-kobo-signature': signature,
        ...options.headers
    };

    const res = await fetch(`${KOBO_API_URL}${endpoint}`, {
        ...options,
        headers,
        body: bodyString ? bodyString : undefined
    });

    if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`Kobo API error ${res.status}: ${errorText}`);
    }

    return res.json();
}
