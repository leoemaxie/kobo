import { env } from '$env/dynamic/private';

const UNSEND_API_URL = env.UNSEND_API_URL || 'https://api.unosend.co/v1/emails';

export interface EmailPayload {
	to: string | string[];
	from?: string;
	subject: string;
	html: string;
	text?: string;
}

export async function sendEmail(payload: EmailPayload) {
	const apiKey = env.UNSEND_API_KEY;
	if (!apiKey) {
		console.warn('[Unsend] UNSEND_API_KEY is not set. Email sending skipped.');
		return { success: false, error: 'Missing API Key' };
	}

	const domain = env.KOBO_DOMAIN || 'triumph.edu';
	const from = payload.from || `Triumph Academy Support <support@${domain}>`;
	const to = Array.isArray(payload.to) ? payload.to : [payload.to];

	try {
		const response = await fetch(UNSEND_API_URL, {
			method: 'POST',
			headers: {
				'Authorization': `Bearer ${apiKey}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				to,
				from,
				subject: payload.subject,
				html: payload.html,
				text: payload.text
			})
		});

		if (!response.ok) {
			const errorText = await response.text();
			console.error(`[Unsend] Error sending email (Status ${response.status}):`, errorText);
			throw new Error('Email sending failed');
		}

		return await response.json();
	} catch (error) {
		console.error('[Unsend] Failed to execute email request:', error);
		throw error;
	}
}
