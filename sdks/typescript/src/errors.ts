import { KoboErrorBody } from "./types.js";

export const KoboErrorCode = {
  IDENTITY_NOT_FOUND: "identity_not_found",
  INVALID_TRANSITION: "invalid_transition",
  DUPLICATE_EXTERNAL_REFERENCE: "duplicate_external_reference",
  INVALID_REQUEST: "invalid_request",
  INVALID_ID: "invalid_id",
  INVALID_QUERY: "invalid_query",
  NOT_FOUND: "not_found",
  INTERNAL_ERROR: "internal_error",
  METHOD_NOT_ALLOWED: "method_not_allowed",
} as const;

export type KoboErrorCodeType = typeof KoboErrorCode[keyof typeof KoboErrorCode];

/**
 * KoboError is thrown when the Kobo API returns a non-2xx response.
 * Always branch on `code`, not `message`; the message is not stable.
 */
export class KoboError extends Error {
  /** HTTP status code */
  readonly status: number;
  /** Stable machine-readable error code, e.g. KoboErrorCode.IDENTITY_NOT_FOUND */
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
