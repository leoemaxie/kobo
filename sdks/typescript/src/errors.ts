import { KoboErrorBody } from "./types.js";

/**
 * KoboError is thrown when the Kobo API returns a non-2xx response.
 * Always branch on `code`, not `message`; the message is not stable.
 */
export class KoboError extends Error {
  /** HTTP status code */
  readonly status: number;
  /** Stable machine-readable error code, e.g. "identity_not_found" */
  readonly code: string;
  /** Optional structured detail fields */
  readonly details?: Record<string, unknown>;

  constructor(status: number, body: KoboErrorBody) {
    super(`kobo [${body.code}]: ${body.message}`);
    this.name = "KoboError";
    this.status = status;
    this.code = body.code;
    if (body.details !== undefined) {
      this.details = body.details;
    }
  }
}
