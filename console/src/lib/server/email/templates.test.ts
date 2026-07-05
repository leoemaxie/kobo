import { describe, it, expect } from 'vitest';
import {
	verificationEmailTemplate,
	keyRotationAlertTemplate,
	billingNoticeTemplate,
	passwordResetTemplate
} from './templates';

describe('Email Templates', () => {
	describe('verificationEmailTemplate', () => {
		it('should generate an email containing the token and default base URL', () => {
			const html = verificationEmailTemplate('test-token-123');
			expect(html).toContain('test-token-123');
			expect(html).toContain('https://console.kobo.dev/auth/verify-email?token=test-token-123');
			expect(html).toContain('Verify your email');
		});

		it('should generate an email with a custom base URL', () => {
			const html = verificationEmailTemplate('token-456', 'https://custom.kobo.dev');
			expect(html).toContain('https://custom.kobo.dev/auth/verify-email?token=token-456');
		});
	});

	describe('keyRotationAlertTemplate', () => {
		it('should contain the environment and keyId', () => {
			const html = keyRotationAlertTemplate('production', 'key_12345');
			expect(html).toContain('production');
			expect(html).toContain('key_12345');
			expect(html).toContain('API Key Rotated');
		});
	});

	describe('billingNoticeTemplate', () => {
		it('should contain the period, amount, and default invoiceUrl', () => {
			const html = billingNoticeTemplate('May 2026', '$100.00');
			expect(html).toContain('May 2026');
			expect(html).toContain('$100.00');
			expect(html).toContain('https://console.kobo.dev/dashboard/billing');
			expect(html).toContain('New Invoice Available');
		});

		it('should support a custom invoiceUrl', () => {
			const html = billingNoticeTemplate('May 2026', '$100.00', 'https://custom-billing.url');
			expect(html).toContain('https://custom-billing.url');
		});
	});

	describe('passwordResetTemplate', () => {
		it('should generate an email containing the token and default base URL', () => {
			const html = passwordResetTemplate('reset-token-xyz');
			expect(html).toContain('reset-token-xyz');
			expect(html).toContain('https://console.kobo.dev/auth/reset-password?token=reset-token-xyz');
			expect(html).toContain('Reset Your Password');
		});

		it('should generate an email with a custom base URL', () => {
			const html = passwordResetTemplate('reset-token-abc', 'https://custom.kobo.dev');
			expect(html).toContain('https://custom.kobo.dev/auth/reset-password?token=reset-token-abc');
		});
	});
});
