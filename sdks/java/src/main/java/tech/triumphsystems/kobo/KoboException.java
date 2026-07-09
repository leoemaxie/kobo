package tech.triumphsystems.kobo;

import java.io.Serial;
import java.util.Collections;
import java.util.Map;

/**
 * Thrown when the Kobo API returns a non-2xx HTTP response.
 *
 * <p>Always branch on {@link #getCode()}, not {@link #getMessage()};
 * the message string is not part of the stable API contract.
 *
 * <pre>{@code
 * try {
 *     kobo.identities().get("unknown-id");
 * } catch (KoboException e) {
 *     if (ErrorCode.IDENTITY_NOT_FOUND.equals(e.getCode())) {
 *         // handle 404
 *     }
 *     throw e;
 * }
 * }</pre>
 */
public final class KoboException extends RuntimeException {

    @Serial
    private static final long serialVersionUID = 1L;

    /** HTTP status code (e.g. 400, 401, 404, 409). */
    private final int httpStatus;

    /**
     * Stable machine-readable error code.
     * Examples: {@link ErrorCode#IDENTITY_NOT_FOUND}, {@link ErrorCode#INVALID_TRANSITION},
     * {@link ErrorCode#DUPLICATE_EXTERNAL_REFERENCE}.
     */
    private final String code;

    /** Arbitrary extra fields returned in the {@code details} field, never null. */
    private final Map<String, Object> details;

    KoboException(int httpStatus, String code, String message, Map<String, Object> details) {
        super("kobo [" + code + "]: " + message);
        this.httpStatus = httpStatus;
        this.code = code;
        this.details = details == null ? Collections.emptyMap() : Collections.unmodifiableMap(details);
    }

    /** @return HTTP status code */
    public int getHttpStatus() { return httpStatus; }

    /** @return Stable machine-readable error code */
    public String getCode() { return code; }

    /** @return Immutable map of extra detail fields; empty if none. */
    public Map<String, Object> getDetails() { return details; }
}
