package tech.triumphsystems.kobo;

/**
 * Constants for Kobo API error codes.
 * These codes are stable and machine-readable, returned in the {@link KoboException#getCode()} field.
 */
public final class ErrorCode {

    private ErrorCode() {
    }

    public static final String IDENTITY_NOT_FOUND = "identity_not_found";
    public static final String INVALID_TRANSITION = "invalid_transition";
    public static final String DUPLICATE_EXTERNAL_REFERENCE = "duplicate_external_reference";
    public static final String INVALID_REQUEST = "invalid_request";
    public static final String INVALID_ID = "invalid_id";
    public static final String INVALID_QUERY = "invalid_query";
    public static final String NOT_FOUND = "not_found";
    public static final String INTERNAL_ERROR = "internal_error";
    public static final String METHOD_NOT_ALLOWED = "method_not_allowed";
}
