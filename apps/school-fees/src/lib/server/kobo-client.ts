import { KoboClient } from '@kobo/sdk';
import { env } from '$env/dynamic/private';

let memoryApiKey: string | null = null;
let memoryApiSecret: string | null = null;

export function setKoboCredentials(key: string, secret: string) {
  memoryApiKey = key;
  memoryApiSecret = secret;
}

function getKoboClient(): KoboClient {
  const apiKey = memoryApiKey || (env.KOBO_API_KEY as string);
  const apiSecret = memoryApiSecret || (env.KOBO_API_SECRET as string);

  const isProd = process.env.NODE_ENV === 'production';

  return isProd
    ? new KoboClient(apiKey, apiSecret)
    : new KoboClient(apiKey, apiSecret, { baseUrl: 'http://localhost:8080/v1' });
}

// Proxy to dynamically resolve the client using the latest credentials
// without needing to rewrite all the routes!
export const kobo = new Proxy({} as KoboClient, {
  get(target, prop, receiver) {
    const client = getKoboClient();
    const value = Reflect.get(client, prop, receiver);
    return typeof value === 'function' ? value.bind(client) : value;
  },
});
