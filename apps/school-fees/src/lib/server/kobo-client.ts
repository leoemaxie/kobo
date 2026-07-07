import { env } from '$env/dynamic/private';

let memoryApiKey: string | null = null;
let memoryApiSecret: string | null = null;

export function setKoboCredentials(key: string, secret: string) {
    memoryApiKey = key;
    memoryApiSecret = secret;
}

export async function koboFetch(endpoint: string, options: RequestInit = {}) {
    const KOBO_API_URL = env.KOBO_API_URL || 'https://api.kobo.dev/v1';
    const apiKey = memoryApiKey || env.KOBO_API_KEY as string;
    const apiSecret = memoryApiSecret || env.KOBO_API_SECRET as string;

    const timestamp = Math.floor(Date.now() / 1000).toString();
    const bodyString = options.body ? (typeof options.body === 'string' ? options.body : JSON.stringify(options.body)) : '';
    
    // Create signature payload: timestamp + "." + body
    const payload = `${timestamp}.${bodyString}`;
    
    // Generate HMAC SHA-256 signature using Web Crypto API
    const encoder = new TextEncoder();
    const key = await crypto.subtle.importKey(
        'raw', 
        encoder.encode(apiSecret), 
        { name: 'HMAC', hash: 'SHA-256' }, 
        false, 
        ['sign']
    );
    
    const signatureBuffer = await crypto.subtle.sign('HMAC', key, encoder.encode(payload));
    const signature = Array.from(new Uint8Array(signatureBuffer))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');

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
