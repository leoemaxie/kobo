import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { EmailService } from './index';
import * as unsend from './unsend';

// Mock the unsend module
vi.mock('./unsend', () => ({
  sendEmail: vi.fn(),
}));

describe('EmailService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should send a password reset email', async () => {
    const sendEmailMock = vi.mocked(unsend.sendEmail);
    sendEmailMock.mockResolvedValue({ success: true } as any);

    await EmailService.sendPasswordResetEmail('test@example.com', 'reset-token-123');

    expect(sendEmailMock).toHaveBeenCalledTimes(1);
    expect(sendEmailMock).toHaveBeenCalledWith({
      to: 'test@example.com',
      subject: 'Reset your Kobo Console password',
      html: expect.stringContaining('reset-token-123'),
    });
  });

  it('should send a verification email', async () => {
    const sendEmailMock = vi.mocked(unsend.sendEmail);
    sendEmailMock.mockResolvedValue({ success: true } as any);

    await EmailService.sendVerificationEmail('test@example.com', 'verify-token-456');

    expect(sendEmailMock).toHaveBeenCalledTimes(1);
    expect(sendEmailMock).toHaveBeenCalledWith({
      to: 'test@example.com',
      subject: 'Verify your Kobo Console account',
      html: expect.stringContaining('verify-token-456'),
    });
  });

  it('should send a key rotation alert email', async () => {
    const sendEmailMock = vi.mocked(unsend.sendEmail);
    sendEmailMock.mockResolvedValue({ success: true } as any);

    await EmailService.sendKeyRotationAlert('test@example.com', 'sandbox', 'key_abc123');

    expect(sendEmailMock).toHaveBeenCalledTimes(1);
    expect(sendEmailMock).toHaveBeenCalledWith({
      to: 'test@example.com',
      subject: '[Kobo] API Key Rotated - SANDBOX',
      html: expect.stringContaining('key_abc123'),
    });
  });

  it('should send a billing notice email', async () => {
    const sendEmailMock = vi.mocked(unsend.sendEmail);
    sendEmailMock.mockResolvedValue({ success: true } as any);

    await EmailService.sendBillingNotice('test@example.com', 'June 2026', '$150.00');

    expect(sendEmailMock).toHaveBeenCalledTimes(1);
    expect(sendEmailMock).toHaveBeenCalledWith({
      to: 'test@example.com',
      subject: '[Kobo] New Invoice Available - June 2026',
      html: expect.stringContaining('$150.00'),
    });
  });
});
