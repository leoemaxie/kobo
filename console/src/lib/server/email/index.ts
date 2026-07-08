import { env } from '$env/dynamic/private';
import { sendEmail } from './unsend';
import { verificationEmailTemplate, keyRotationAlertTemplate, billingNoticeTemplate, passwordResetTemplate, invitationEmailTemplate } from './templates';

const getDomain = () => env.KOBO_DOMAIN || 'kobo.dev';

export const EmailService = {
	async sendPasswordResetEmail(to: string, token: string, baseUrl?: string) {
		return sendEmail({
			to,
			from: `Kobo Security <security@${getDomain()}>`,
			subject: 'Reset your Kobo Console password',
			html: passwordResetTemplate(token, baseUrl)
		});
	},

	async sendVerificationEmail(to: string, token: string, baseUrl?: string) {
		return sendEmail({
			to,
			from: `Kobo Accounts <accounts@${getDomain()}>`,
			subject: 'Verify your Kobo Console account',
			html: verificationEmailTemplate(token, baseUrl)
		});
	},

	async sendKeyRotationAlert(to: string, environment: string, keyId: string) {
		return sendEmail({
			to,
			from: `Kobo Security <security@${getDomain()}>`,
			subject: `[Kobo] API Key Rotated - ${environment.toUpperCase()}`,
			html: keyRotationAlertTemplate(environment, keyId)
		});
	},

	async sendBillingNotice(to: string, period: string, amount: string, invoiceUrl?: string) {
		return sendEmail({
			to,
			from: `Kobo Billing <billing@${getDomain()}>`,
			subject: `[Kobo] New Invoice Available - ${period}`,
			html: billingNoticeTemplate(period, amount, invoiceUrl)
		});
	},

	async sendInvitationEmail(to: string, role: string, workspaceName: string, token: string, baseUrl?: string) {
		return sendEmail({
			to,
			from: `Kobo Invitations <invites@${getDomain()}>`,
			subject: `You have been invited to join ${workspaceName} on Kobo`,
			html: invitationEmailTemplate(role, workspaceName, token, baseUrl)
		});
	}
};
