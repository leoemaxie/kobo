export function verificationEmailTemplate(
  token: string,
  baseUrl: string = "https://console.kobo.dev",
): string {
  const verifyUrl = `${baseUrl}/auth/verify-email?token=${token}`;
  return `
		<div style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #f5f5f5; background-color: #0a0a0a;">
			<div style="background-color: #111; border: 1px solid #222; border-radius: 12px; padding: 40px; text-align: center;">
				<div style="margin-bottom: 24px;">
					<span style="font-size: 28px; font-weight: 800; color: #c6f135; letter-spacing: -1px;">k.</span>
				</div>
				<h1 style="font-size: 24px; font-weight: 700; color: #f5f5f5; margin: 0 0 16px 0; letter-spacing: -0.5px;">Verify your email</h1>
				<p style="font-size: 15px; color: #888; line-height: 1.6; margin: 0 0 32px 0;">
					Welcome to Kobo. Please verify your email address to complete your registration and unlock your environment.
				</p>
				<a href="${verifyUrl}" style="display: inline-block; padding: 14px 28px; background-color: #c6f135; color: #080808; text-decoration: none; font-weight: 700; border-radius: 6px; font-size: 14px; letter-spacing: -0.2px;">
					Verify Email Address
				</a>
				<p style="font-size: 13px; color: #555; margin-top: 40px; padding-top: 24px; border-top: 1px solid #222;">
					If you did not request this, please ignore this email. The link will expire in 24 hours.
				</p>
			</div>
		</div>
	`;
}

export function keyRotationAlertTemplate(
  environment: string,
  keyId: string,
): string {
  return `
		<div style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #f5f5f5; background-color: #0a0a0a;">
			<div style="background-color: #111; border: 1px solid #222; border-radius: 12px; padding: 40px;">
				<div style="margin-bottom: 24px; text-align: center;">
					<span style="font-size: 28px; font-weight: 800; color: #c6f135; letter-spacing: -1px;">k.</span>
				</div>
				<h1 style="font-size: 20px; font-weight: 700; color: #f5f5f5; margin: 0 0 16px 0; letter-spacing: -0.5px; text-align: center;">API Key Rotated</h1>
				<p style="font-size: 15px; color: #888; line-height: 1.6; margin: 0 0 24px 0; text-align: center;">
					Your Kobo API key for the <strong>${environment}</strong> environment has been successfully rotated.
				</p>
				<div style="background: #0f0f0f; border: 1px solid #1e1e1e; padding: 16px; border-radius: 8px; margin: 0 0 24px 0;">
					<p style="font-family: monospace; font-size: 13px; color: #aaa; margin: 0;">
						<strong style="color: #666; text-transform: uppercase; letter-spacing: 0.5px; margin-right: 8px;">Key ID</strong> ${keyId}
					</p>
				</div>
				<p style="font-size: 14px; color: #666; line-height: 1.5; margin: 0; padding: 16px; background: rgba(239, 68, 68, 0.08); border: 1px solid rgba(239, 68, 68, 0.2); border-radius: 8px; color: #f87171;">
					The previous secret key is no longer valid. Please ensure your infrastructure has been updated with the new secret.
				</p>
			</div>
		</div>
	`;
}

export function billingNoticeTemplate(
  period: string,
  amount: string,
  invoiceUrl: string = "https://console.kobo.dev/dashboard/billing",
): string {
  return `
		<div style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #f5f5f5; background-color: #0a0a0a;">
			<div style="background-color: #111; border: 1px solid #222; border-radius: 12px; padding: 40px; text-align: center;">
				<div style="margin-bottom: 24px;">
					<span style="font-size: 28px; font-weight: 800; color: #c6f135; letter-spacing: -1px;">k.</span>
				</div>
				<h1 style="font-size: 20px; font-weight: 700; color: #f5f5f5; margin: 0 0 16px 0; letter-spacing: -0.5px;">New Invoice Available</h1>
				<p style="font-size: 15px; color: #888; line-height: 1.6; margin: 0 0 32px 0;">
					Your invoice for the period of <strong style="color: #aaa;">${period}</strong> is now available in your Kobo Console.
				</p>
				<div style="background: #0f0f0f; border: 1px solid #1e1e1e; padding: 24px; border-radius: 8px; margin: 0 0 32px 0;">
					<p style="font-size: 12px; color: #666; text-transform: uppercase; letter-spacing: 0.5px; margin: 0 0 8px 0; font-weight: 600;">Amount Due</p>
					<p style="font-size: 32px; font-weight: 700; color: #c6f135; margin: 0; letter-spacing: -1px;">${amount}</p>
				</div>
				<a href="${invoiceUrl}" style="display: inline-block; padding: 14px 28px; background-color: #c6f135; color: #080808; text-decoration: none; font-weight: 700; border-radius: 6px; font-size: 14px; letter-spacing: -0.2px;">
					View & Pay Invoice
				</a>
			</div>
		</div>
	`;
}

export function passwordResetTemplate(
  token: string,
  baseUrl: string = "https://console.kobo.dev",
): string {
  const resetUrl = `${baseUrl}/auth/reset-password?token=${token}`;
  return `
		<div style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #f5f5f5; background-color: #0a0a0a;">
			<div style="background-color: #111; border: 1px solid #222; border-radius: 12px; padding: 40px; text-align: center;">
				<div style="margin-bottom: 24px;">
					<span style="font-size: 28px; font-weight: 800; color: #c6f135; letter-spacing: -1px;">k.</span>
				</div>
				<h1 style="font-size: 20px; font-weight: 700; color: #f5f5f5; margin: 0 0 16px 0; letter-spacing: -0.5px;">Reset Your Password</h1>
				<p style="font-size: 15px; color: #888; line-height: 1.6; margin: 0 0 32px 0;">
					We received a request to reset the password for your Kobo Console account.
				</p>
				<a href="${resetUrl}" style="display: inline-block; padding: 14px 28px; background-color: #c6f135; color: #080808; text-decoration: none; font-weight: 700; border-radius: 6px; font-size: 14px; letter-spacing: -0.2px;">
					Reset Password
				</a>
				<p style="font-size: 13px; color: #555; margin-top: 40px; padding-top: 24px; border-top: 1px solid #222;">
					If you did not request this password reset, please ignore this email or contact support if you have concerns.
				</p>
			</div>
		</div>
	`;
}

export function invitationEmailTemplate(
  role: string,
  workspaceName: string,
  token: string,
  baseUrl: string = "https://console.kobo.dev",
): string {
  const inviteUrl = `${baseUrl}/auth/accept-invite?token=${token}`;
  return `
		<div style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #f5f5f5; background-color: #0a0a0a;">
			<div style="background-color: #111; border: 1px solid #222; border-radius: 12px; padding: 40px; text-align: center;">
				<div style="margin-bottom: 24px;">
					<span style="font-size: 28px; font-weight: 800; color: #c6f135; letter-spacing: -1px;">k.</span>
				</div>
				<h1 style="font-size: 20px; font-weight: 700; color: #f5f5f5; margin: 0 0 16px 0; letter-spacing: -0.5px;">You're invited!</h1>
				<p style="font-size: 15px; color: #888; line-height: 1.6; margin: 0 0 32px 0;">
					You have been invited to join the <strong style="color: #f5f5f5;">${workspaceName}</strong> workspace on Kobo as a <strong style="color: #f5f5f5; text-transform: capitalize;">${role}</strong>.
				</p>
				<a href="${inviteUrl}" style="display: inline-block; padding: 14px 28px; background-color: #c6f135; color: #080808; text-decoration: none; font-weight: 700; border-radius: 6px; font-size: 14px; letter-spacing: -0.2px;">
					Accept Invitation
				</a>
				<p style="font-size: 13px; color: #555; margin-top: 40px; padding-top: 24px; border-top: 1px solid #222;">
					If you did not expect this invitation, you can ignore this email. This link expires in 7 days.
				</p>
			</div>
		</div>
	`;
}
