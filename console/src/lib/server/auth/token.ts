import { createHash } from 'crypto';

export function generateToken(): string {
	const bytes = new Uint8Array(32);
	crypto.getRandomValues(bytes);
	return Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
}

export function hashToken(token: string): string {
	return createHash('sha256').update(token).digest('hex');
}
