import { db } from '$lib/server/db';
import { invitations } from '$lib/server/db/schema';
import type { InferInsertModel } from 'drizzle-orm';

// In a real app, you would send an email here using an email provider like Resend or SendGrid.
export async function createInvitation(
	integratorId: string,
	invitedBy: string,
	email: string,
	role: 'member' | 'owner' | 'superadmin'
) {
	const randomBytesId = new Uint8Array(16);
	crypto.getRandomValues(randomBytesId);
	const token = Array.from(randomBytesId).map((b) => b.toString(16).padStart(2, '0')).join('');

	// Expires in 7 days
	const expiresAt = new Date(Date.now() + 1000 * 60 * 60 * 24 * 7);

	const invitation: InferInsertModel<typeof invitations> = {
		id: token,
		integratorId,
		invitedBy,
		email,
		role,
		expiresAt
	};

	await db.insert(invitations).values(invitation);

	// Mock sending email
	console.log(`Sending invitation email to ${email} with token ${token}`);

	return token;
}
