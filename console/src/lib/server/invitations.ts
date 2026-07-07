import { db } from '$lib/server/db';
import { invitations } from '$lib/server/db/schema';
import type { InferInsertModel } from 'drizzle-orm';
import { env } from '$env/dynamic/private';
import { sendEmail } from '$lib/server/email/unsend';

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

	const domain = env.KOBO_DOMAIN || 'kobo.dev';
	const baseUrl = env.PUBLIC_URL || 'https://console.kobo.dev';
	const inviteUrl = `${baseUrl}/auth/register?token=${token}`;

	await sendEmail({
		to: email,
		from: `Kobo Invitations <invites@${domain}>`,
		subject: 'You have been invited to Kobo',
		html: `
		<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2>You're invited!</h2>
			<p>You have been invited to join a Kobo workspace as a <strong>${role}</strong>.</p>
			<p><a href="${inviteUrl}" style="display: inline-block; padding: 10px 20px; background-color: #65a30d; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Accept Invitation</a></p>
			<p style="font-size: 12px; color: #666; margin-top: 20px;">If you did not expect this invitation, you can ignore this email. This link expires in 7 days.</p>
		</div>`
	});

	return token;
}
