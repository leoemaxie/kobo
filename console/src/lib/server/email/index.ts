import { sendEmail } from './unsend';
import { verificationEmailTemplate, keyRotationAlertTemplate, billingNoticeTemplate, passwordResetTemplate } from './templates';

export const EmailService = {
	async sendPasswordResetEmail(to: string, token: string, baseUrl?: string) {
		return sendEmail({
			to,
			subject: 'Reset your Kobo Console password',
			html: passwordResetTemplate(token, baseUrl)
		});
	},

	async sendVerificationEmail(to: string, token: string, baseUrl?: string) {
		return sendEmail({
			to,
			subject: 'Verify your Kobo Console account',
			html: verificationEmailTemplate(token, baseUrl)
		});
	},

	async sendKeyRotationAlert(to: string, environment: string, keyId: string) {
		return sendEmail({
			to,
			subject: `[Kobo] API Key Rotated - ${environment.toUpperCase()}`,
			html: keyRotationAlertTemplate(environment, keyId)
		});
	},

	async sendBillingNotice(to: string, period: string, amount: string, invoiceUrl?: string) {
		return sendEmail({
			to,
			subject: `[Kobo] New Invoice Available - ${period}`,
			html: billingNoticeTemplate(period, amount, invoiceUrl)
		});
	}
};
