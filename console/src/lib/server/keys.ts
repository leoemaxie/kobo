import * as argon2 from 'argon2';

export async function generateKeyPair(env: 'sandbox' | 'production') {
	const prefix = env === 'production' ? 'sk_live_' : 'sk_test_';
	const randomBytesId = new Uint8Array(8);
	crypto.getRandomValues(randomBytesId);
	const randomId = Array.from(randomBytesId).map(b => b.toString(16).padStart(2, '0')).join('');
	const keyId = `${prefix}${randomId}`;

	const randomBytesSecret = new Uint8Array(32);
	crypto.getRandomValues(randomBytesSecret);
	const plainSecret = Array.from(randomBytesSecret).map(b => b.toString(16).padStart(2, '0')).join('');

	const secretHash = await argon2.hash(plainSecret);

	return { keyId, plainSecret, secretHash };
}
