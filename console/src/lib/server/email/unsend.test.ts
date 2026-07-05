import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { sendEmail } from './unsend';

// Mock the environment variables
vi.mock('$env/dynamic/private', () => ({
	env: {
		UNSEND_API_KEY: 'test_api_key',
		UNSEND_API_URL: 'https://test.api/emails'
	}
}));

describe('unsend.ts', () => {
	let fetchMock: any;

	beforeEach(() => {
		fetchMock = vi.fn();
		global.fetch = fetchMock;
		
		// Reset mock console.warn and console.error
		vi.spyOn(console, 'warn').mockImplementation(() => {});
		vi.spyOn(console, 'error').mockImplementation(() => {});
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('should send an email successfully', async () => {
		fetchMock.mockResolvedValue({
			ok: true,
			json: async () => ({ id: 'email_123' })
		});

		const result = await sendEmail({
			to: 'test@example.com',
			subject: 'Test Subject',
			html: '<p>Test HTML</p>'
		});

		expect(result).toEqual({ id: 'email_123' });
		expect(fetchMock).toHaveBeenCalledTimes(1);
		
		const fetchCallArgs = fetchMock.mock.calls[0];
		expect(fetchCallArgs[0]).toBe('https://test.api/emails');
		expect(fetchCallArgs[1].method).toBe('POST');
		expect(fetchCallArgs[1].headers).toEqual({
			'Authorization': 'Bearer test_api_key',
			'Content-Type': 'application/json'
		});
		
		const body = JSON.parse(fetchCallArgs[1].body);
		expect(body).toEqual({
			to: ['test@example.com'],
			from: 'Kobo Support <support@kobo.dev>',
			subject: 'Test Subject',
			html: '<p>Test HTML</p>'
		});
	});

	it('should support array of recipients', async () => {
		fetchMock.mockResolvedValue({
			ok: true,
			json: async () => ({ id: 'email_456' })
		});

		await sendEmail({
			to: ['user1@example.com', 'user2@example.com'],
			subject: 'Test Subject',
			html: '<p>Test HTML</p>'
		});

		const body = JSON.parse(fetchMock.mock.calls[0][1].body);
		expect(body.to).toEqual(['user1@example.com', 'user2@example.com']);
	});

	it('should throw an error if the API request fails', async () => {
		fetchMock.mockResolvedValue({
			ok: false,
			status: 400,
			text: async () => 'Bad Request'
		});

		await expect(sendEmail({
			to: 'test@example.com',
			subject: 'Test Subject',
			html: '<p>Test</p>'
		})).rejects.toThrow('Email sending failed');

		expect(console.error).toHaveBeenCalledWith(
			'[Unsend] Error sending email (Status 400):',
			'Bad Request'
		);
	});

	it('should handle fetch throwing an exception', async () => {
		const networkError = new Error('Network error');
		fetchMock.mockRejectedValue(networkError);

		await expect(sendEmail({
			to: 'test@example.com',
			subject: 'Test',
			html: '<p>Test</p>'
		})).rejects.toThrow('Network error');

		expect(console.error).toHaveBeenCalledWith(
			'[Unsend] Failed to execute email request:',
			networkError
		);
	});
});
