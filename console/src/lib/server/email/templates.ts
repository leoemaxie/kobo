export function verificationEmailTemplate(token: string, baseUrl: string = 'https://console.kobo.dev'): string {
	const verifyUrl = `${baseUrl}/auth/verify-email?token=${token}`;
	return `
		<div style="font-family: 'Inter', sans-serif; max-width: 600px; margin: 0 auto; padding: 24px; color: #080808; background-color: #fcfcfc; border: 1px solid #eee; border-radius: 8px;">
			<h1 style="font-size: 24px; font-weight: 700; color: #080808; margin-top: 0;">Verify your email</h1>
			<p style="font-size: 15px; color: #555; line-height: 1.6;">
				Welcome to Kobo. Please verify your email address to complete your registration and unlock your sandbox environment.
			</p>
			<div style="margin: 32px 0;">
				<a href="${verifyUrl}" style="display: inline-block; padding: 12px 24px; background-color: #080808; color: #C0FF00; text-decoration: none; font-weight: 600; border-radius: 6px; font-size: 14px;">
					Verify Email Address
				</a>
			</div>
			<p style="font-size: 13px; color: #888; margin-top: 32px; padding-top: 16px; border-top: 1px solid #eee;">
				If you did not request this, please ignore this email.
			</p>
		</div>
	`;
}

export function keyRotationAlertTemplate(environment: string, keyId: string): string {
	return `
		<div style="font-family: 'Inter', sans-serif; max-width: 600px; margin: 0 auto; padding: 24px; color: #080808; border: 1px solid #eee; border-radius: 8px;">
			<h1 style="font-size: 20px; font-weight: 700; margin-top: 0;">API Key Rotated</h1>
			<p style="font-size: 15px; color: #555; line-height: 1.6;">
				Your Kobo API key for the <strong>${environment}</strong> environment has been successfully rotated.
			</p>
			<div style="background: #f5f5f5; padding: 12px; border-radius: 6px; margin: 20px 0;">
				<p style="font-family: monospace; font-size: 13px; color: #333; margin: 0;">
					<strong>Key ID:</strong> ${keyId}
				</p>
			</div>
			<p style="font-size: 14px; color: #444; line-height: 1.5;">
				The previous secret key is no longer valid. Please ensure your infrastructure has been updated with the new secret.
			</p>
		</div>
	`;
}

export function billingNoticeTemplate(period: string, amount: string, invoiceUrl: string = 'https://console.kobo.dev/dashboard/billing'): string {
	return `
		<div style="font-family: 'Inter', sans-serif; max-width: 600px; margin: 0 auto; padding: 24px; color: #080808; border: 1px solid #eee; border-radius: 8px;">
			<h1 style="font-size: 20px; font-weight: 700; margin-top: 0;">New Invoice Available</h1>
			<p style="font-size: 15px; color: #555; line-height: 1.6;">
				Your invoice for the period of <strong>${period}</strong> is now available in your Kobo Console.
			</p>
			<div style="margin: 24px 0;">
				<p style="font-size: 13px; color: #888; text-transform: uppercase; letter-spacing: 0.05em; margin: 0 0 4px 0;">Amount Due</p>
				<p style="font-size: 24px; font-weight: 600; color: #080808; margin: 0;">${amount}</p>
			</div>
			<div style="margin: 32px 0 16px 0;">
				<a href="${invoiceUrl}" style="display: inline-block; padding: 10px 20px; background-color: #080808; color: #C0FF00; text-decoration: none; font-weight: 600; border-radius: 6px; font-size: 13px;">
					View & Pay Invoice
				</a>
			</div>
		</div>
	`;
}

export function passwordResetTemplate(token: string, baseUrl: string = 'https://console.kobo.dev'): string {
	const resetUrl = `${baseUrl}/auth/reset-password?token=${token}`;
	return `
		<div style="font-family: 'Inter', sans-serif; max-width: 600px; margin: 0 auto; padding: 24px; color: #080808; border: 1px solid #eee; border-radius: 8px;">
			<h1 style="font-size: 20px; font-weight: 700; margin-top: 0;">Reset Your Password</h1>
			<p style="font-size: 15px; color: #555; line-height: 1.6;">
				We received a request to reset the password for your Kobo Console account.
			</p>
			<div style="margin: 32px 0 16px 0;">
				<a href="${resetUrl}" style="display: inline-block; padding: 10px 20px; background-color: #080808; color: #C0FF00; text-decoration: none; font-weight: 600; border-radius: 6px; font-size: 13px;">
					Reset Password
				</a>
			</div>
			<p style="font-size: 13px; color: #888; margin-top: 32px; padding-top: 16px; border-top: 1px solid #eee;">
				If you did not request this password reset, please ignore this email or contact support if you have concerns.
			</p>
		</div>
	`;
}

